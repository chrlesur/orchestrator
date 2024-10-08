package job

import (
	"time"

	"orchestrator/internal/model"
)

// NewJob crée une nouvelle instance de Job
func NewJob(name string, command string) *model.Job {
	return &model.Job{
		ID:        generateJobID(),
		Name:      name,
		Command:   command,
		Status:    model.JobStatusPending,
		CreatedAt: time.Now(),
	}
}

// generateJobID génère un nouvel ID de job commençant par "J"
func generateJobID() string {
	return "J" + time.Now().Format("20060102150405")
}
