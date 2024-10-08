package job

import (
	"fmt"
	"orchestrator/internal/logging"
	"sync"
)

type Manager struct {
	jobs map[string]*Job
	mu   sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		jobs: make(map[string]*Job),
	}
}

func (m *Manager) AddJob(job *Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.jobs[job.ID]; exists {
		return fmt.Errorf("job with ID %s already exists", job.ID)
	}

	m.jobs[job.ID] = job
	logging.InfoLogger.Printf("Added job: %s (%s)", job.Name, job.ID)
	return nil
}

func (m *Manager) GetJob(id string) (*Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	job, exists := m.jobs[id]
	if !exists {
		return nil, fmt.Errorf("job with ID %s not found", id)
	}

	return job, nil
}

func (m *Manager) ListJobs() []*Job {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jobs := make([]*Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobs = append(jobs, job)
	}

	return jobs
}

func (m *Manager) RunJob(id string) error {
	job, err := m.GetJob(id)
	if err != nil {
		return err
	}

	return job.RunWithRetries()
}
