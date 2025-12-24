package handlers

import (
	"encoding/json"
	"net/http"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/usecases"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type IterationHandlers struct {
	iterationUseCase *usecases.IterationUseCase
}

func NewIterationHandlers(iterationUseCase *usecases.IterationUseCase) *IterationHandlers {
	return &IterationHandlers{
		iterationUseCase: iterationUseCase,
	}
}

type CreateIterationRequest struct {
	ProjectID   uuid.UUID `json:"project_id"`
	Number      int       `json:"number"`
	Description string    `json:"description"`
	StartAt     string    `json:"start_at"`
	EndAt       string    `json:"end_at"`
}

// GetAll handles GET /iterations
// @Summary Get all iterations
// @Description Get all iterations for a specific project
// @Tags iterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id query string true "Project ID" format(uuid)
// @Success 200 {array} models.Iteration "List of iterations"
// @Failure 400 {string} string "Invalid project_id"
// @Failure 500 {string} string "Failed to retrieve iterations"
// @Router /iterations [get]
func (h *IterationHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr == "" {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	iterations, err := h.iterationUseCase.GetAll(ctx, projectID)
	if err != nil {
		http.Error(w, "Failed to retrieve iterations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(iterations)
}

// GetByID handles GET /iterations/{id}
// @Summary Get iteration by ID
// @Description Get a specific iteration by its ID
// @Tags iterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Iteration ID" format(uuid)
// @Success 200 {object} models.Iteration "Iteration details"
// @Failure 400 {string} string "Invalid iteration ID"
// @Failure 404 {string} string "Iteration not found"
// @Router /iterations/{id} [get]
func (h *IterationHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid iteration ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	iteration, err := h.iterationUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Iteration not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(iteration)
}

// Create handles POST /iterations
// @Summary Create a new iteration
// @Description Create a new iteration for a project
// @Tags iterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param iteration body CreateIterationRequest true "Iteration data"
// @Success 201 {object} map[string]interface{} "Iteration created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create iteration"
// @Router /iterations [post]
func (h *IterationHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateIterationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startAt, err := parseTime(req.StartAt)
	if err != nil {
		http.Error(w, "Invalid start_at format", http.StatusBadRequest)
		return
	}

	endAt, err := parseTime(req.EndAt)
	if err != nil {
		http.Error(w, "Invalid end_at format", http.StatusBadRequest)
		return
	}

	newIteration := models.Iteration{
		ProjectID:   req.ProjectID,
		Number:      req.Number,
		Description: req.Description,
		StartAt:     startAt,
		EndAt:       endAt,
	}

	ctx := r.Context()
	iterationID, err := h.iterationUseCase.Create(ctx, newIteration)
	if err != nil {
		http.Error(w, "Failed to create iteration", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":          iterationID,
		"project_id":  req.ProjectID,
		"number":      req.Number,
		"description": req.Description,
		"start_at":    req.StartAt,
		"end_at":      req.EndAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Delete handles DELETE /iterations/{id}
// @Summary Delete iteration
// @Description Delete an iteration by its ID
// @Tags iterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Iteration ID" format(uuid)
// @Success 204 "Iteration deleted successfully"
// @Failure 400 {string} string "Invalid iteration ID"
// @Failure 500 {string} string "Failed to delete iteration"
// @Router /iterations/{id} [delete]
func (h *IterationHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid iteration ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.iterationUseCase.Delete(ctx, id)
	if err != nil {
		http.Error(w, "Failed to delete iteration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetIterationAnalysis handles GET /iterations/{id}/analysis
// @Summary Get iteration indicator analysis
// @Description Get detailed analysis of iteration indicators with data points for graphing
// @Tags iterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Iteration ID" format(uuid)
// @Success 200 {object} models.IterationAnalysisResponse "Iteration analysis with indicator data points"
// @Failure 400 {string} string "Invalid iteration ID"
// @Failure 404 {string} string "Iteration not found"
// @Failure 500 {string} string "Failed to retrieve analysis"
// @Router /iterations/{id}/analysis [get]
func (h *IterationHandlers) GetIterationAnalysis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid iteration ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	analysis, err := h.iterationUseCase.GetIterationAnalysis(ctx, id)
	if err != nil {
		http.Error(w, "Failed to retrieve analysis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}
