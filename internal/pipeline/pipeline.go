package pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/chrlesur/orchestrator/internal/job"
	"github.com/chrlesur/orchestrator/internal/models"
	"github.com/chrlesur/orchestrator/pkg/logger"
)

func NewPipeline(id, name string, jobs []*models.Job, scheduledAt time.Time) *models.Pipeline {
	return &models.Pipeline{
		ID:          id,
		Name:        name,
		Jobs:        jobs,
		Status:      models.PipelineStatusPending,
		Context:     make(map[string]interface{}),
		ScheduledAt: scheduledAt,
	}
}

func Execute(p *models.Pipeline, ctx context.Context) error {
	p.Status = models.PipelineStatusRunning
	p.StartTime = time.Now()
	defer func() { p.EndTime = time.Now() }()

	for _, j := range p.Jobs {
		err := job.Execute(j, ctx) // Utilisez job.Execute au lieu de j.Execute
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
