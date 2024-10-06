package job

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chrlesur/orchestrator/internal/db"
	"github.com/chrlesur/orchestrator/internal/models"
	"github.com/chrlesur/orchestrator/internal/plugin"
	"github.com/chrlesur/orchestrator/pkg/logger"
	"github.com/chrlesur/orchestrator/pkg/utils"
)

type Manager struct {
	jobs           map[string]*models.Job
	jobQueue       chan *models.Job
	mu             sync.Mutex
	wg             sync.WaitGroup
	store          *db.Store
	defaultTimeout time.Duration
	maxRetries     int
	pluginManager  *plugin.PluginManager
}

func NewManager(workerCount int, store *db.Store, defaultTimeout time.Duration, maxRetries int, pluginManager *plugin.PluginManager) *Manager {
	m := &Manager{
		jobs:           make(map[string]*models.Job),
		jobQueue:       make(chan *models.Job, 100),
		store:          store,
		defaultTimeout: defaultTimeout,
		maxRetries:     maxRetries,
		pluginManager:  pluginManager,
	}

	// Charger les jobs existants depuis la base de données
	jobs, err := store.GetAllJobs()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load jobs from database: %v", err))
	} else {
		for _, job := range jobs {
			m.jobs[job.ID] = job
		}
	}

	for i := 0; i < workerCount; i++ {
		go m.worker()
	}

	return m
}

func (m *Manager) CreateJob(command string, args []string, pluginName string) (*models.Job, error) {
	id := utils.GenerateID(8)
	job := &models.Job{
		ID:         id,
		Command:    command,
		Args:       args,
		PluginName: pluginName,
		Status:     models.JobStatusPending,
		Timeout:    m.defaultTimeout,
		MaxRetries: m.maxRetries,
	}

	err := m.AddJob(job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (m *Manager) AddJob(job *models.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.jobs[job.ID]; exists {
		return fmt.Errorf("job with ID %s already exists", job.ID)
	}

	m.jobs[job.ID] = job
	err := m.store.SaveJob(job)
	if err != nil {
		return fmt.Errorf("failed to save job to database: %v", err)
	}

	m.jobQueue <- job
	m.wg.Add(1)

	return nil
}

func (m *Manager) GetJob(id string) (*models.Job, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.jobs[id]
	if !exists {
		// Si le job n'est pas en mémoire, essayons de le récupérer depuis la base de données
		dbJob, err := m.store.GetJob(id)
		if err != nil {
			return nil, fmt.Errorf("job with ID %s not found", id)
		}
		m.jobs[id] = dbJob
		return dbJob, nil
	}

	return job, nil
}

func (m *Manager) GetJobs() []*models.Job {
	m.mu.Lock()
	defer m.mu.Unlock()

	jobs := make([]*models.Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (m *Manager) worker() {
	for job := range m.jobQueue {
		logger.Info(fmt.Sprintf("Starting job %s", job.ID))
		start := time.Now()
		var err error
		if job.PluginName != "" {
			err = m.executePluginJob(job)
		} else {
			err = Execute(job, context.Background())
		}
		duration := time.Since(start)
		if err != nil {
			logger.Error(fmt.Sprintf("Job %s failed after %s: %v", job.ID, utils.FormatDuration(duration), err))
		} else {
			logger.Info(fmt.Sprintf("Job %s completed successfully in %s", job.ID, utils.FormatDuration(duration)))
		}
		m.store.SaveJob(job) // Sauvegarder l'état final du job
		m.wg.Done()
	}
}

func (m *Manager) executePluginJob(job *models.Job) error {
	args := make(map[string]interface{})
	for i, arg := range job.Args {
		args[fmt.Sprintf("arg%d", i)] = arg
	}

	result, err := m.pluginManager.ExecutePlugin(job.PluginName, args)
	if err != nil {
		job.Status = models.JobStatusFailed
		job.Error = err
		return err
	}

	job.Result = fmt.Sprintf("%v", result)
	job.Status = models.JobStatusCompleted
	return nil
}

func (m *Manager) Wait() {
	m.wg.Wait()
}

func (m *Manager) Shutdown() {
	close(m.jobQueue)
	m.Wait()
}

func (m *Manager) UpdateJob(job *models.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.jobs[job.ID]; !exists {
		return fmt.Errorf("job with ID %s not found", job.ID)
	}

	m.jobs[job.ID] = job
	return m.store.SaveJob(job)
}
