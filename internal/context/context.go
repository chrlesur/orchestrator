package context

import (
	"encoding/json"
	"time"

	"orchestrator/internal/model"
)

// Context représente le contexte d'exécution d'un job ou d'un pipeline
type Context struct {
	ID        string          `json:"id"`
	Content   json.RawMessage `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// PipelineContext représente le contexte global d'un pipeline
type PipelineContext struct {
	PipelineID   string             `json:"pipeline_id"`
	JobContexts  map[string]Context `json:"job_contexts"`
	FinalContext Context            `json:"final_context"`
}

// NewContext crée un nouveau contexte
func NewContext(content interface{}) (*model.Context, error) {
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	return &model.Context{
		ID:        generateContextID(),
		Content:   jsonContent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// NewPipelineContext crée un nouveau contexte de pipeline
func NewPipelineContext(pipelineID string) *PipelineContext {
	return &PipelineContext{
		PipelineID:  pipelineID,
		JobContexts: make(map[string]Context),
	}
}

// AddJobContext ajoute le contexte d'un job au contexte du pipeline
func (pc *PipelineContext) AddJobContext(jobID string, jobContext Context) {
	pc.JobContexts[jobID] = jobContext
}

// SetFinalContext définit le contexte final du pipeline
func (pc *PipelineContext) SetFinalContext(finalContext Context) {
	pc.FinalContext = finalContext
}

// UpdatePipelineContext met à jour le contexte du pipeline après l'exécution d'un job
func UpdatePipelineContext(pc *PipelineContext, job job.Job, jobContext Context) {
	pc.AddJobContext(job.ID, jobContext)

	// Si c'est le dernier job du pipeline, on définit le contexte final
	p, _ := pipeline.GetPipeline(pc.PipelineID)
	if p != nil && p.Jobs[len(p.Jobs)-1].ID == job.ID {
		pc.SetFinalContext(jobContext)
	}
}

// generateContextID génère un nouvel ID de contexte
func generateContextID() string {
	return "C" + time.Now().Format("20060102150405")
}
