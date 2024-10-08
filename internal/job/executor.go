package job

import (
    "context"
    "os/exec"
    "time"
    "orchestrator/internal/logging"
)

func (j *Job) Execute() error {
    j.Status = StatusRunning
    logging.InfoLogger.Printf("Starting job execution: %s (%s)", j.Name, j.ID)

    ctx, cancel := context.WithTimeout(context.Background(), j.Timeout)
    defer cancel()

    cmd := exec.CommandContext(ctx, "sh", "-c", j.Command)
    cmd.Dir = j.WorkDir

    output, err := cmd.CombinedOutput()
    j.Context.SetContent(string(output))

    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            j.Status = StatusFailed
            logging.ErrorLogger.Printf("Job timed out: %s (%s)", j.Name, j.ID)
            return fmt.Errorf("job timed out after %v", j.Timeout)
        }

        j.Status = StatusFailed
        logging.ErrorLogger.Printf("Job failed: %s (%s) - %v", j.Name, j.ID, err)
        return err
    }

    j.Status = StatusCompleted
    logging.InfoLogger.Printf("Job completed successfully: %s (%s)", j.Name, j.ID)
    return nil
}

func (j *Job) RunWithRetries() error {
    var lastErr error

    for j.CurrentRetry <= j.MaxRetries {
        err := j.Execute()
        if err == nil {
            return nil
        }

        lastErr = err
        j.CurrentRetry++

        if j.CurrentRetry <= j.MaxRetries {
            logging.InfoLogger.Printf("Retrying job: %s (%s) - Attempt %d/%d", j.Name, j.ID, j.CurrentRetry, j.MaxRetries)
            time.Sleep(time.Second * time.Duration(j.CurrentRetry)) // Exponential backoff
        }
    }

    logging.ErrorLogger.Printf("Job failed after %d retries: %s (%s)", j.MaxRetries, j.Name, j.ID)
    return lastErr
}