package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chrlesur/orchestrator/internal/db"
	"github.com/chrlesur/orchestrator/internal/job"
	"github.com/chrlesur/orchestrator/internal/models"
	"github.com/chrlesur/orchestrator/internal/plugin"
	"github.com/chrlesur/orchestrator/pkg/logger"
)

type Manager struct {
	pipelines     map[string]*models.Pipeline
	pipelineQueue chan *models.Pipeline
	mu            sync.Mutex
	wg            sync.WaitGroup
	store         *db.Store
	pluginManager *plugin.PluginManager
}

func NewManager(workerCount int, store *db.Store, pluginManager *plugin.PluginManager) *Manager {
	m := &Manager{
		pipelines:     make(map[string]*models.Pipeline),
		pipelineQueue: make(chan *models.Pipeline, 100),
		store:         store,
		pluginManager: pluginManager,
	}

	// Charger les pipelines existants depuis la base de données
	pipelines, err := store.GetAllPipelines()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load pipelines from database: %v", err))
	} else {
		for _, pipeline := range pipelines {
			m.pipelines[pipeline.ID] = pipeline
		}
	}

	for i := 0; i < workerCount; i++ {
		go m.worker()
	}

	go m.scheduler()

	return m
}

func (m *Manager) AddPipeline(pipeline *models.Pipeline) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.pipelines[pipeline.ID]; exists {
		return fmt.Errorf("pipeline with ID %s already exists", pipeline.ID)
	}

	m.pipelines[pipeline.ID] = pipeline
	err := m.store.SavePipeline(pipeline)
	if err != nil {
		return fmt.Errorf("failed to save pipeline to database: %v", err)
	}

	return nil
}

func (m *Manager) GetPipeline(id string) (*models.Pipeline, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pipeline, exists := m.pipelines[id]
	if !exists {
		// Si le pipeline n'est pas en mémoire, essayons de le récupérer depuis la base de données
		dbPipeline, err := m.store.GetPipeline(id)
		if err != nil {
			return nil, fmt.Errorf("pipeline with ID %s not found", id)
		}
		m.pipelines[id] = dbPipeline
		return dbPipeline, nil
	}

	return pipeline, nil
}

func (m *Manager) UpdatePipeline(pipeline *models.Pipeline) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	existingPipeline, exists := m.pipelines[pipeline.ID]
	if !exists {
		return fmt.Errorf("pipeline with ID %s not found", pipeline.ID)
	}

	// Mettre à jour les champs du pipeline existant
	existingPipeline.Name = pipeline.Name
	existingPipeline.Jobs = pipeline.Jobs

	// Sauvegarder les modifications dans la base de données
	err := m.store.SavePipeline(existingPipeline)
	if err != nil {
		return fmt.Errorf("failed to save updated pipeline to database: %v", err)
	}

	logger.Info(fmt.Sprintf("Pipeline %s updated", pipeline.ID))
	return nil
}

func (m *Manager) DeletePipeline(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.pipelines[id]; !exists {
		return fmt.Errorf("pipeline with ID %s not found", id)
	}

	delete(m.pipelines, id)
	// Implémentez la méthode DeletePipeline dans le store si nécessaire
	// err := m.store.DeletePipeline(id)
	// if err != nil {
	//     return fmt.Errorf("failed to delete pipeline from database: %v", err)
	// }

	logger.Info(fmt.Sprintf("Pipeline %s deleted", id))
	return nil
}

func (m *Manager) GetPipelines() []*models.Pipeline {
	m.mu.Lock()
	defer m.mu.Unlock()

	pipelines := make([]*models.Pipeline, 0, len(m.pipelines))
	for _, pipeline := range m.pipelines {
		pipelines = append(pipelines, pipeline)
	}
	return pipelines
}

func (m *Manager) worker() {
	for pipeline := range m.pipelineQueue {
		logger.Info(fmt.Sprintf("Starting pipeline %s", pipeline.ID))
		err := m.executePipeline(pipeline)
		if err != nil {
			logger.Error(fmt.Sprintf("Pipeline %s failed: %v", pipeline.ID, err))
		} else {
			logger.Info(fmt.Sprintf("Pipeline %s completed successfully", pipeline.ID))
		}
		m.store.SavePipeline(pipeline) // Sauvegarder l'état final du pipeline
		m.wg.Done()
	}
}

func (m *Manager) executePipeline(p *models.Pipeline) error {
	p.Status = models.PipelineStatusRunning
	p.StartTime = time.Now()
	defer func() { p.EndTime = time.Now() }()

	for _, j := range p.Jobs {
		var err error
		if j.PluginName != "" {
			args := make(map[string]interface{})
			for i, arg := range j.Args {
				args[fmt.Sprintf("arg%d", i)] = arg
			}
			_, err = m.pluginManager.ExecutePlugin(j.PluginName, args)
		} else {
			err = job.Execute(j, context.Background())
		}
		if err != nil {
			p.Status = models.PipelineStatusFailed
			logger.Error(fmt.Sprintf("Pipeline %s failed: job %s encountered an error: %v", p.ID, j.ID, err))
			return err
		}

		// Agréger le contexte du job dans le contexte du pipeline
		p.Context[j.ID] = j.Result
	}

	p.Status = models.PipelineStatusCompleted
	logger.Info(fmt.Sprintf("Pipeline %s completed successfully", p.ID))
	return nil
}

func (m *Manager) scheduler() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		now := time.Now()
		m.mu.Lock()
		for _, pipeline := range m.pipelines {
			if pipeline.Status == models.PipelineStatusPending && now.After(pipeline.ScheduledAt) {
				m.pipelineQueue <- pipeline
				m.wg.Add(1)
				pipeline.Status = models.PipelineStatusRunning
				m.store.SavePipeline(pipeline) // Sauvegarder le changement de statut
			}
		}
		m.mu.Unlock()
	}
}

func (m *Manager) Wait() {
	m.wg.Wait()
}

func (m *Manager) Shutdown() {
	close(m.pipelineQueue)
	m.Wait()
}
