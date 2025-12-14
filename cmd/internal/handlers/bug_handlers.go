package handlers

import (
	"encoding/json"
	"net/http"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/usecases"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BugHandlers struct {
	bugUseCase *usecases.BugUseCase
}

func NewBugHandlers(bugUseCase *usecases.BugUseCase) *BugHandlers {
	return &BugHandlers{
		bugUseCase: bugUseCase,
	}
}

type CreateBugRequest struct {
	TaskID      uuid.UUID  `json:"task_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	Number      int        `json:"number"`
	Description string     `json:"description"`
	Points      int        `json:"points"`
}

// GetAll handles GET /bugs
// @Summary Get all bugs
// @Description Get all bugs for a specific task
// @Tags bugs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id query string true "Task ID" format(uuid)
// @Success 200 {array} models.Bug "List of bugs"
// @Failure 400 {string} string "Invalid task_id"
// @Failure 500 {string} string "Failed to retrieve bugs"
// @Router /bugs [get]
func (h *BugHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
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
	bugs, err := h.bugUseCase.GetAll(ctx, taskID)
	if err != nil {
		http.Error(w, "Failed to retrieve bugs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bugs)
}

// GetByID handles GET /bugs/{id}
// @Summary Get bug by ID
// @Description Get a specific bug by its ID
// @Tags bugs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bug ID" format(uuid)
// @Success 200 {object} models.Bug "Bug details"
// @Failure 400 {string} string "Invalid bug ID"
// @Failure 404 {string} string "Bug not found"
// @Router /bugs/{id} [get]
func (h *BugHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid bug ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	bug, err := h.bugUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Bug not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bug)
}

// Create handles POST /bugs
// @Summary Create a new bug
// @Description Create a new bug for a task
// @Tags bugs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bug body CreateBugRequest true "Bug data"
// @Success 201 {object} map[string]interface{} "Bug created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create bug"
// @Router /bugs [post]
func (h *BugHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateBugRequest
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

	newBug := models.Bug{
		TaskID:      req.TaskID,
		Assignee:    assignee,
		Number:      req.Number,
		Description: req.Description,
		Points:      points,
	}

	ctx := r.Context()
	bugID, err := h.bugUseCase.Create(ctx, newBug)
	if err != nil {
		http.Error(w, "Failed to create bug", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":          bugID,
		"task_id":     req.TaskID,
		"number":      req.Number,
		"description": req.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
