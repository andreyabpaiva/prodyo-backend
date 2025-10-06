package handlers

import (
	"encoding/json"
	"net/http"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/usecases"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProjectHandlers struct {
	projectUseCase *usecases.ProjectUseCase
}

func NewProjectHandlers(projectUseCase *usecases.ProjectUseCase) *ProjectHandlers {
	return &ProjectHandlers{
		projectUseCase: projectUseCase,
	}
}

type CreateProjectRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateProjectRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// GetAllProjects handles GET /projects
func (h *ProjectHandlers) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projects, err := h.projectUseCase.GetAll(ctx)
	if err != nil {
		http.Error(w, "Failed to retrieve projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// GetProjectByID handles GET /projects/{id}
func (h *ProjectHandlers) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	project, err := h.projectUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// CreateProject handles POST /projects
func (h *ProjectHandlers) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	newProject := models.Project{
		Name:  req.Name,
		Email: req.Email,
	}

	ctx := r.Context()
	projectID, err := h.projectUseCase.Add(ctx, newProject)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":    projectID,
		"name":  req.Name,
		"email": req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateProject handles PUT /projects/{id}
func (h *ProjectHandlers) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	updatedProject := models.Project{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}

	ctx := r.Context()
	err = h.projectUseCase.Update(ctx, updatedProject)
	if err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProject)
}

// DeleteProject handles DELETE /projects/{id}
func (h *ProjectHandlers) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.projectUseCase.Delete(ctx, id)
	if err != nil {
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
