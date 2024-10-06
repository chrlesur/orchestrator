package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chrlesur/orchestrator/internal/job"
	"github.com/chrlesur/orchestrator/internal/models"
	"github.com/chrlesur/orchestrator/internal/pipeline"
	"github.com/chrlesur/orchestrator/internal/plugin"
	"github.com/gorilla/mux"
)

type Server struct {
	jobManager      *job.Manager
	pipelineManager *pipeline.Manager
	pluginManager   *plugin.PluginManager
	router          *mux.Router
}

func NewServer(jobManager *job.Manager, pipelineManager *pipeline.Manager, pluginManager *plugin.PluginManager) *Server {
	s := &Server{
		jobManager:      jobManager,
		pipelineManager: pipelineManager,
		pluginManager:   pluginManager,
		router:          mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/jobs", authMiddleware(s.handleGetJobs)).Methods("GET")
	s.router.HandleFunc("/jobs", authMiddleware(s.handleCreateJob)).Methods("POST")
	s.router.HandleFunc("/jobs/{id}", authMiddleware(s.handleGetJob)).Methods("GET")
	s.router.HandleFunc("/pipelines", authMiddleware(s.handleGetPipelines)).Methods("GET")
	s.router.HandleFunc("/pipelines", authMiddleware(s.handleCreatePipeline)).Methods("POST")
	s.router.HandleFunc("/pipelines/{id}", authMiddleware(s.handleGetPipeline)).Methods("GET")
	s.router.HandleFunc("/plugins", authMiddleware(s.handleGetPlugins)).Methods("GET")
	s.router.HandleFunc("/plugins/{name}/execute", authMiddleware(s.handleExecutePlugin)).Methods("POST")
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) handleGetJobs(w http.ResponseWriter, r *http.Request) {
	jobs := s.jobManager.GetJobs()
	respondJSON(w, http.StatusOK, jobs)
}

func (s *Server) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var jobReq struct {
		ID      string   `json:"id"`
		Command string   `json:"command"`
		Args    []string `json:"args"`
	}

	if err := json.NewDecoder(r.Body).Decode(&jobReq); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newJob, err := s.jobManager.CreateJob(jobReq.Command, jobReq.Args, "")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, newJob)
}

func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]

	job, err := s.jobManager.GetJob(jobID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Job not found")
		return
	}

	respondJSON(w, http.StatusOK, job)
}

func (s *Server) handleGetPipelines(w http.ResponseWriter, r *http.Request) {
	pipelines := s.pipelineManager.GetPipelines()
	respondJSON(w, http.StatusOK, pipelines)
}

func (s *Server) handleCreatePipeline(w http.ResponseWriter, r *http.Request) {
	var pipelineReq struct {
		ID     string   `json:"id"`
		Name   string   `json:"name"`
		JobIDs []string `json:"job_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&pipelineReq); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	jobs := make([]*models.Job, 0, len(pipelineReq.JobIDs))
	for _, jobID := range pipelineReq.JobIDs {
		job, err := s.jobManager.GetJob(jobID)
		if err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Job %s not found", jobID))
			return
		}
		jobs = append(jobs, job)
	}

	newPipeline := &models.Pipeline{
		ID:          pipelineReq.ID,
		Name:        pipelineReq.Name,
		Jobs:        jobs,
		Status:      models.PipelineStatusPending,
		ScheduledAt: time.Now().Add(1 * time.Minute),
	}
	if err := s.pipelineManager.AddPipeline(newPipeline); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, newPipeline)
}

func (s *Server) handleGetPipeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pipelineID := vars["id"]

	pipeline, err := s.pipelineManager.GetPipeline(pipelineID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Pipeline not found")
		return
	}

	respondJSON(w, http.StatusOK, pipeline)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

func (s *Server) handleStartJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]

	job, err := s.jobManager.GetJob(jobID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Job not found")
		return
	}

	err = s.jobManager.AddJob(job)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to start job: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "Job started"})
}

func (s *Server) handleStopJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]

	job, err := s.jobManager.GetJob(jobID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Job not found")
		return
	}

	// Implement job stopping logic here
	// For now, we'll just update the status
	job.Status = models.JobStatusFailed
	err = s.jobManager.UpdateJob(job)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to stop job: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "Job stopped"})
}

func (s *Server) handleUpdatePipeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pipelineID := vars["id"]

	var pipelineUpdate struct {
		Name   string   `json:"name"`
		JobIDs []string `json:"job_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&pipelineUpdate); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	pipeline, err := s.pipelineManager.GetPipeline(pipelineID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Pipeline not found")
		return
	}

	jobs := make([]*models.Job, 0, len(pipelineUpdate.JobIDs))
	for _, jobID := range pipelineUpdate.JobIDs {
		job, err := s.jobManager.GetJob(jobID)
		if err != nil {
			respondError(w, http.StatusBadRequest, fmt.Sprintf("Job %s not found", jobID))
			return
		}
		jobs = append(jobs, job)
	}

	pipeline.Name = pipelineUpdate.Name
	pipeline.Jobs = jobs

	err = s.pipelineManager.UpdatePipeline(pipeline)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update pipeline: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, pipeline)
}

func (s *Server) handleGetPlugins(w http.ResponseWriter, r *http.Request) {
	plugins := s.pluginManager.GetLoadedPlugins()
	respondJSON(w, http.StatusOK, plugins)
}

func (s *Server) handleExecutePlugin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pluginName := vars["name"]

	var pluginArgs map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&pluginArgs); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	result, err := s.pluginManager.ExecutePlugin(pluginName, pluginArgs)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to execute plugin: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, result)
}
