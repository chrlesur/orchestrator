package server

import (
	"errors"
	"orchestrator/internal/job"
	"sync"

	"orchestrator/internal/context"
	"orchestrator/internal/pipeline"
	"orchestrator/internal/storage"
)

var (
	pipelines      = make(map[string]*pipeline.Pipeline)
	pipelinesMutex sync.RWMutex
)

// CreatePipeline crée un nouveau pipeline
func CreatePipeline(name string, jobIDs []string) (*pipeline.Pipeline, error) {
	jobs, err := getJobsByIDs(jobIDs)
	if err != nil {
		return nil, err
	}

	newPipeline := pipeline.NewPipeline(name, jobs)
	err = storage.SavePipeline(newPipeline)
	if err != nil {
		return nil, err
	}
	return newPipeline, nil
}

// ListPipelines retourne la liste de tous les pipelines
func ListPipelines() ([]*pipeline.Pipeline, error) {
	// Cette fonction nécessite une implémentation spécifique dans storage/boltdb.go
	return storage.ListPipelines()
}

// GetPipeline récupère les détails d'un pipeline spécifique
func GetPipeline(id string) (*pipeline.Pipeline, error) {
	return storage.GetPipeline(id)
}

// RunPipeline exécute un pipeline spécifique
func RunPipeline(id string) error {
	p, err := GetPipeline(id)
	if err != nil {
		return err
	}

	p.Status = pipeline.StatusRunning
	storage.SavePipeline(p)

	pipelineContext := context.NewPipelineContext(id)

	for _, job := range p.Jobs {
		jobContext, err := RunJob(job.ID)
		if err != nil {
			p.Status = pipeline.StatusFailed
			storage.SavePipeline(p)
			return err
		}

		context.UpdatePipelineContext(pipelineContext, job, *jobContext)
	}

	p.Status = pipeline.StatusCompleted
	storage.SavePipeline(p)

	// Sauvegarde du contexte du pipeline
	storage.SaveContext(&pipelineContext.FinalContext)

	return nil
}

// getJobsByIDs récupère les jobs correspondant aux IDs fournis
func getJobsByIDs(jobIDs []string) ([]job.Job, error) {
	// À implémenter : récupérer les jobs à partir de leur ID
	// Cette fonction dépendra de votre implémentation de la gestion des jobs
	return nil, errors.New("not implemented")
}
