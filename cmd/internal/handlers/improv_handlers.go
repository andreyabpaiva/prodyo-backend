package handlers

import (
	"encoding/json"
	"net/http"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/usecases"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ImprovHandlers struct {
	improvUseCase *usecases.ImprovUseCase
}

func NewImprovHandlers(improvUseCase *usecases.ImprovUseCase) *ImprovHandlers {
	return &ImprovHandlers{
		improvUseCase: improvUseCase,
	}
}

type CreateImprovRequest struct {
	TaskID      uuid.UUID  `json:"task_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	Number      int        `json:"number"`
	Description string     `json:"description"`
	Points      int        `json:"points"`
}

// GetAll handles GET /improvements
// @Summary Get all improvements
// @Description Get all improvements for a specific task
// @Tags improvements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id query string true "Task ID" format(uuid)
// @Success 200 {array} models.Improv "List of improvements"
// @Failure 400 {string} string "Invalid task_id"
// @Failure 500 {string} string "Failed to retrieve improvements"
// @Router /improvements [get]
func (h *ImprovHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Query().Get("task_id")
	if taskIDStr == "" {
		http.Error(w, "task_id is required", http.StatusBadRequest)
		return
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	improvements, err := h.improvUseCase.GetAll(ctx, taskID)
	if err != nil {
		http.Error(w, "Failed to retrieve improvements", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(improvements)
}

// GetByID handles GET /improvements/{id}
// @Summary Get improvement by ID
// @Description Get a specific improvement by its ID
// @Tags improvements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Improvement ID" format(uuid)
// @Success 200 {object} models.Improv "Improvement details"
// @Failure 400 {string} string "Invalid improvement ID"
// @Failure 404 {string} string "Improvement not found"
// @Router /improvements/{id} [get]
func (h *ImprovHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid improvement ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	improv, err := h.improvUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Improvement not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(improv)
}

// Create handles POST /improvements
// @Summary Create a new improvement
// @Description Create a new improvement for a task
// @Tags improvements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param improvement body CreateImprovRequest true "Improvement data"
// @Success 201 {object} map[string]interface{} "Improvement created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create improvement"
// @Router /improvements [post]
func (h *ImprovHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateImprovRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var assignee models.User
	if req.AssigneeID != nil {
		assignee.ID = *req.AssigneeID
	}

	points := req.Points
	if points == 0 {
		points = 1
	}

	newImprov := models.Improv{
		TaskID:      req.TaskID,
		Assignee:    assignee,
		Number:      req.Number,
		Description: req.Description,
		Points:      points,
	}

	ctx := r.Context()
	improvID, err := h.improvUseCase.Create(ctx, newImprov)
	if err != nil {
		http.Error(w, "Failed to create improvement", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":          improvID,
		"task_id":     req.TaskID,
		"number":      req.Number,
		"description": req.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
