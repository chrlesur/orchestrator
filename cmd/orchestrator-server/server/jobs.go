package server

import (
	"net/http"
	"orchestrator/internal/context"
	"orchestrator/internal/job"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) createJob(c *gin.Context) {
	var jobRequest struct {
		Name       string        `json:"name" binding:"required"`
		Command    string        `json:"command" binding:"required"`
		WorkDir    string        `json:"workDir" binding:"required"`
		Timeout    time.Duration `json:"timeout" binding:"required"`
		MaxRetries int           `json:"maxRetries" binding:"required"`
	}

	if err := c.ShouldBindJSON(&jobRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newJob := job.NewJob(
		job.GenerateID(),
		jobRequest.Name,
		jobRequest.Command,
		jobRequest.WorkDir,
		jobRequest.Timeout,
		jobRequest.MaxRetries,
	)

	if err := s.jobManager.AddJob(newJob); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newJob)
}

func (s *Server) listJobs(c *gin.Context) {
	jobs := s.jobManager.ListJobs()
	c.JSON(http.StatusOK, jobs)
}

func (s *Server) getJob(c *gin.Context) {
	id := c.Param("id")
	job, err := s.jobManager.GetJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, job)
}

func RunJob(id string) (*context.Context, error) {
	j, err := GetJob(id)
	if err != nil {
		return nil, err
	}

	j.Status = job.StatusRunning

	// Exécution du job ici...
	// Pour l'exemple, on va simplement créer un contexte factice
	result := map[string]string{"output": "Job executed successfully"}

	jobContext, err := context.NewContext(result)
	if err != nil {
		j.Status = job.StatusFailed
		return nil, err
	}

	j.Status = job.StatusCompleted

	// Ici, vous devriez sauvegarder le contexte du job dans la base de données
	// saveJobContext(jobContext)

	return jobContext, nil
}
