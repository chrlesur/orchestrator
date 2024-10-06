package job

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/chrlesur/orchestrator/internal/models"
	"github.com/chrlesur/orchestrator/pkg/logger"
)

func NewJob(id, command string, args []string, timeout time.Duration, maxRetries int) *models.Job {
	return &models.Job{
		ID:         id,
		Command:    command,
		Args:       args,
		Timeout:    timeout,
		MaxRetries: maxRetries,
		Status:     models.JobStatusPending,
	}
}

func Execute(j *models.Job, ctx context.Context) error {
	j.Status = models.JobStatusRunning
	j.StartTime = time.Now()
	defer func() { j.EndTime = time.Now() }()

	for j.RetryCount <= j.MaxRetries {
		err := run(j, ctx)
		if err == nil {
			j.Status = models.JobStatusCompleted
			return nil
		}

		j.RetryCount++
		j.Error = err
		logger.Warning(fmt.Sprintf("Job %s failed (attempt %d/%d): %v", j.ID, j.RetryCount, j.MaxRetries+1, err))

		if j.RetryCount > j.MaxRetries {
			break
		}

		// Attente exponentielle entre les tentatives
		time.Sleep(time.Second * time.Duration(1<<uint(j.RetryCount)))
	}

	j.Status = models.JobStatusFailed
	return fmt.Errorf("job %s failed after %d attempts: %v", j.ID, j.RetryCount, j.Error)
}

func run(j *models.Job, ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, j.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, j.Command, j.Args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution failed: %v, output: %s", err, string(output))
	}

	j.Result = string(output)
	return nil
}
