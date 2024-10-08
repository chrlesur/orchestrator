package model

import (
	"encoding/json"
	"time"
)

// JobStatus représente l'état actuel d'un job
type JobStatus string

const (
	JobStatusPending   JobStatus = "PENDING"
	JobStatusRunning   JobStatus = "RUNNING"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed    JobStatus = "FAILED"
)

// Job représente un job à exécuter
type Job struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	Status    JobStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// PipelineStatus représente l'état actuel d'un pipeline
type PipelineStatus string

const (
	PipelineStatusPending   PipelineStatus = "PENDING"
	PipelineStatusRunning   PipelineStatus = "RUNNING"
	PipelineStatusCompleted PipelineStatus = "COMPLETED"
	PipelineStatusFailed    PipelineStatus = "FAILED"
)

// Pipeline représente un pipeline d'exécution de jobs
type Pipeline struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	JobIDs          []string       `json:"job_ids"`
	Status          PipelineStatus `json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	LastExecutionAt time.Time      `json:"last_execution_at"`
}

// Context représente le contexte d'exécution d'un job ou d'un pipeline
type Context struct {
	ID        string          `json:"id"`
	Content   json.RawMessage `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
