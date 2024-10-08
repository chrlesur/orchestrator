package pipeline

import (
	"time"

	"orchestrator/internal/model"
)

// NewPipeline crée une nouvelle instance de Pipeline
func NewPipeline(name string, jobIDs []string) *model.Pipeline {
	return &model.Pipeline{
		ID:        generatePipelineID(),
		Name:      name,
		JobIDs:    jobIDs,
		Status:    model.PipelineStatusPending,
		CreatedAt: time.Now(),
	}
}

// generatePipelineID génère un nouvel ID de pipeline commençant par "P"
func generatePipelineID() string {
	return "P" + time.Now().Format("20060102150405")
}
