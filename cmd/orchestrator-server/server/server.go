package server

import (
	"fmt"
	"net/http"
	"orchestrator/internal/config"
	"orchestrator/internal/job"
	"orchestrator/internal/logging"
	"orchestrator/internal/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg        *config.Config
	jobManager *job.Manager
	router     *gin.Engine
}

func Init() error {
	return storage.InitDB("orchestrator.db")
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		cfg:        cfg,
		jobManager: job.NewManager(),
		router:     gin.Default(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.POST("/jobs", s.createJob)
	s.router.GET("/jobs", s.listJobs)
	s.router.GET("/jobs/:id", s.getJob)
	s.router.POST("/jobs/:id/run", s.runJob)
}

func (s *Server) Run() {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	logging.InfoLogger.Printf("Starting server on %s", addr)
	if err := s.router.Run(addr); err != nil {
		logging.ErrorLogger.Printf("Failed to start server: %v", err)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Routes existantes pour les jobs
	r.POST("/jobs", createJob)
	r.GET("/jobs", listJobs)
	r.GET("/jobs/:id", getJob)
	r.POST("/jobs/:id/run", runJob)
	r.GET("/jobs/:id/context", getJobContext)

	// Nouvelles routes pour les pipelines
	r.POST("/pipelines", createPipeline)
	r.GET("/pipelines", listPipelines)
	r.GET("/pipelines/:id", getPipeline)
	r.POST("/pipelines/:id/run", runPipeline)
	r.GET("/pipelines/:id/context", getPipelineContext)

	return r
}

func createPipeline(c *gin.Context) {
	var request struct {
		Name   string   `json:"name" binding:"required"`
		JobIDs []string `json:"job_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPipeline, err := CreatePipeline(request.Name, request.JobIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPipeline)
}

func listPipelines(c *gin.Context) {
	pipelines := ListPipelines()
	c.JSON(http.StatusOK, pipelines)
}

func getPipeline(c *gin.Context) {
	id := c.Param("id")
	pipeline, err := GetPipeline(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}
	c.JSON(http.StatusOK, pipeline)
}

func runPipeline(c *gin.Context) {
	id := c.Param("id")
	err := RunPipeline(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pipeline execution started"})
}

func getPipelineContext(c *gin.Context) {
	id := c.Param("id")
	// Récupérer le contexte du pipeline depuis la base de données
	// pipelineContext := getPipelineContextFromDB(id)

	// Pour l'exemple, on renvoie un contexte factice
	c.JSON(http.StatusOK, gin.H{"message": "Pipeline context for " + id})
}

func getJobContext(c *gin.Context) {
	id := c.Param("id")
	// Récupérer le contexte du job depuis la base de données
	// jobContext := getJobContextFromDB(id)

	// Pour l'exemple, on renvoie un contexte factice
	c.JSON(http.StatusOK, gin.H{"message": "Job context for " + id})
}
