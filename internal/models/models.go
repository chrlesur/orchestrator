package models

import (
	"time"
)

type JobStatus string
type PipelineStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"

	PipelineStatusPending   PipelineStatus = "pending"
	PipelineStatusRunning   PipelineStatus = "running"
	PipelineStatusCompleted PipelineStatus = "completed"
	PipelineStatusFailed    PipelineStatus = "failed"
)

type Job struct {
	ID         string
	Command    string
	Args       []string
	Timeout    time.Duration
	MaxRetries int
	Status     JobStatus
	Result     string
	Error      error
	StartTime  time.Time
	EndTime    time.Time
	RetryCount int
	PluginName string
}

type Pipeline struct {
	ID          string
	Name        string
	Jobs        []*Job
	Status      PipelineStatus
	StartTime   time.Time
	EndTime     time.Time
	Context     map[string]interface{}
	ScheduledAt time.Time
}
