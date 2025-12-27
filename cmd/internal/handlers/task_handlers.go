package handlers

import (
	"encoding/json"
	"net/http"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/usecases"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TaskHandlers struct {
	taskUseCase *usecases.TaskUseCase
}

func NewTaskHandlers(taskUseCase *usecases.TaskUseCase) *TaskHandlers {
	return &TaskHandlers{
		taskUseCase: taskUseCase,
	}
}

type CreateTaskRequest struct {
	IterationID uuid.UUID  `json:"iteration_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	Status      string     `json:"status"`
	Timer       *string    `json:"timer,omitempty"`
	Points      int        `json:"points"`
}

type UpdateTaskRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	Status      string     `json:"status"`
	Timer       *string    `json:"timer,omitempty"`
	Points      int        `json:"points"`
}

type PatchTaskRequest struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Timer       *string    `json:"timer,omitempty"`
	Points      *int       `json:"points,omitempty"`
}

// GetAll handles GET /tasks
// @Summary Get all tasks
// @Description Get all tasks for a specific iteration
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param iteration_id query string true "Iteration ID" format(uuid)
// @Success 200 {array} models.Task "List of tasks"
// @Failure 400 {string} string "Invalid iteration_id"
// @Failure 500 {string} string "Failed to retrieve tasks"
// @Router /tasks [get]
func (h *TaskHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	iterationIDStr := r.URL.Query().Get("iteration_id")
	if iterationIDStr == "" {
		http.Error(w, "iteration_id is required", http.StatusBadRequest)
		return
	}

	iterationID, err := uuid.Parse(iterationIDStr)
	if err != nil {
		http.Error(w, "Invalid iteration_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tasks, err := h.taskUseCase.GetAll(ctx, iterationID)
	if err != nil {
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetByID handles GET /tasks/{id}
// @Summary Get task by ID
// @Description Get a specific task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} models.Task "Task details"
// @Failure 400 {string} string "Invalid task ID"
// @Failure 404 {string} string "Task not found"
// @Router /tasks/{id} [get]
func (h *TaskHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	task, err := h.taskUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Create handles POST /tasks
// @Summary Create a new task
// @Description Create a new task for an iteration
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body CreateTaskRequest true "Task data"
// @Success 201 {object} map[string]interface{} "Task created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create task"
// @Router /tasks [post]
func (h *TaskHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var assignee models.User
	if req.AssigneeID != nil {
		assignee.ID = *req.AssigneeID
	}

	var timer int64
	if req.Timer != nil {
		t, err := parseDuration(*req.Timer)
		if err == nil {
			timer = t
		}
	}

	status := normalizeStatus(req.Status)

	points := req.Points
	if points == 0 {
		points = 1
	}

	newTask := models.Task{
		IterationID: req.IterationID,
		Name:        req.Name,
		Description: req.Description,
		Assignee:    assignee,
		Status:      status,
		Timer:       timer,
		Points:      points,
	}

	ctx := r.Context()
	taskID, err := h.taskUseCase.Create(ctx, newTask)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":           taskID,
		"iteration_id": req.IterationID,
		"name":         req.Name,
		"description":  req.Description,
		"status":       status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Update handles PUT /tasks/{id}
// @Summary Update task
// @Description Update an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID" format(uuid)
// @Param task body UpdateTaskRequest true "Updated task data"
// @Success 200 {object} models.Task "Updated task"
// @Failure 400 {string} string "Invalid task ID or request body"
// @Failure 500 {string} string "Failed to update task"
// @Router /tasks/{id} [put]
func (h *TaskHandlers) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var assignee models.User
	if req.AssigneeID != nil {
		assignee.ID = *req.AssigneeID
	}

	var timer int64
	if req.Timer != nil {
		t, err := parseDuration(*req.Timer)
		if err == nil {
			timer = t
		}
	}

	status := normalizeStatus(req.Status)

	points := req.Points
	if points == 0 {
		points = 1
	}

	updatedTask := models.Task{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Assignee:    assignee,
		Status:      status,
		Timer:       timer,
		Points:      points,
	}

	ctx := r.Context()
	err = h.taskUseCase.Update(ctx, updatedTask)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

// Delete handles DELETE /tasks/{id}
// @Summary Delete task
// @Description Delete a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID" format(uuid)
// @Success 204 "Task deleted successfully"
// @Failure 400 {string} string "Invalid task ID"
// @Failure 500 {string} string "Failed to delete task"
// @Router /tasks/{id} [delete]
func (h *TaskHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.taskUseCase.Delete(ctx, id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Patch handles PATCH /tasks/{id}
// @Summary Partially update task
// @Description Partially update an existing task (only provided fields will be updated)
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID" format(uuid)
// @Param task body PatchTaskRequest true "Partial task data"
// @Success 200 {object} models.Task "Updated task"
// @Failure 400 {string} string "Invalid task ID or request body"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Failed to update task"
// @Router /tasks/{id} [patch]
func (h *TaskHandlers) Patch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	existingTask, err := h.taskUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	var req PatchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		existingTask.Name = *req.Name
	}

	if req.Description != nil {
		existingTask.Description = *req.Description
	}

	if req.AssigneeID != nil {
		existingTask.Assignee.ID = *req.AssigneeID
	}

	if req.Status != nil {
		existingTask.Status = normalizeStatus(*req.Status)
	}

	if req.Timer != nil {
		t, err := parseDuration(*req.Timer)
		if err == nil {
			existingTask.Timer = t
		}
	}

	if req.Points != nil {
		points := *req.Points
		if points == 0 {
			points = 1
		}
		existingTask.Points = points
	}

	err = h.taskUseCase.Update(ctx, existingTask)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	updatedTask, err := h.taskUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Failed to retrieve updated task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}
