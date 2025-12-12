package handlers

import (
	"encoding/json"
	"net/http"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/usecases"
	"strconv"

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
	Name        string                   `json:"name" validate:"required"`
	Description string                   `json:"description"`
	Color       string                   `json:"color"`
	ProdRange   models.ProductivityRange `json:"prod_range"`
	MemberIDs   []string                 `json:"member_ids"`
}

type UpdateProjectRequest struct {
	Name        string                   `json:"name" validate:"required"`
	Description string                   `json:"description"`
	Color       string                   `json:"color"`
	ProdRange   models.ProductivityRange `json:"prod_range"`
	MemberIDs   []string                 `json:"member_ids"`
}

// GetAllProjects handles GET /projects
// @Summary Get all projects
// @Description Get a paginated list of all projects with their members
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Success 200 {object} map[string]interface{} "Projects with pagination"
// @Failure 500 {string} string "Internal server error"
// @Router /projects [get]
func (h *ProjectHandlers) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pagination := models.PaginationRequest{
		Page:     1,
		PageSize: 20,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pagination.Page = p
		}
	}

	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			pagination.PageSize = ps
		}
	}

	projects, paginationResp, err := h.projectUseCase.GetAll(ctx, pagination)
	if err != nil {
		http.Error(w, "Failed to retrieve projects", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       projects,
		"pagination": paginationResp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProjectByID handles GET /projects/{id}
// @Summary Get project by ID
// @Description Get a specific project by its ID with all members
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID" format(uuid)
// @Success 200 {object} models.Project "Project details"
// @Failure 400 {string} string "Invalid project ID"
// @Failure 404 {string} string "Project not found"
// @Router /projects/{id} [get]
func (h *ProjectHandlers) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	project, iterationCount, err := h.projectUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":               project.ID,
		"name":             project.Name,
		"description":      project.Description,
		"color":            project.Color,
		"prod_range":       project.ProdRange,
		"members":          project.Members,
		"created_at":       project.CreatedAt,
		"updated_at":       project.UpdatedAt,
		"iterations_count": iterationCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateProject handles POST /projects
// @Summary Create a new project
// @Description Create a new project with members and productivity range
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project body CreateProjectRequest true "Project data"
// @Success 201 {object} map[string]interface{} "Created project"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create project"
// @Router /projects [post]
func (h *ProjectHandlers) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	var members []models.User
	for _, memberID := range req.MemberIDs {
		if id, err := uuid.Parse(memberID); err == nil {
			members = append(members, models.User{ID: id})
		}
	}

	newProject := models.Project{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		ProdRange:   req.ProdRange,
		Members:     members,
	}

	ctx := r.Context()
	projectID, err := h.projectUseCase.Add(ctx, newProject)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":          projectID,
		"name":        req.Name,
		"description": req.Description,
		"color":       req.Color,
		"prod_range":  req.ProdRange,
		"member_ids":  req.MemberIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateProject handles PUT /projects/{id}
// @Summary Update project
// @Description Update an existing project with new data
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID" format(uuid)
// @Param project body UpdateProjectRequest true "Updated project data"
// @Success 200 {object} models.Project "Updated project"
// @Failure 400 {string} string "Invalid project ID or request body"
// @Failure 500 {string} string "Failed to update project"
// @Router /projects/{id} [put]
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

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Convert member IDs to User objects
	var members []models.User
	for _, memberID := range req.MemberIDs {
		if id, err := uuid.Parse(memberID); err == nil {
			members = append(members, models.User{ID: id})
		}
	}

	updatedProject := models.Project{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		ProdRange:   req.ProdRange,
		Members:     members,
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
// @Summary Delete project
// @Description Delete a project by its ID
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID" format(uuid)
// @Success 204 "Project deleted successfully"
// @Failure 400 {string} string "Invalid project ID"
// @Failure 500 {string} string "Failed to delete project"
// @Router /projects/{id} [delete]
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

// GetProjectsByMemberID handles GET /projects/member/{userId}
// @Summary Get projects by member ID
// @Description Get all projects where the specified user is a member
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID" format(uuid)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Success 200 {object} map[string]interface{} "Projects with pagination"
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Internal server error"
// @Router /projects/member/{userId} [get]
func (h *ProjectHandlers) GetProjectsByMemberID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	pagination := models.PaginationRequest{
		Page:     1,
		PageSize: 20,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pagination.Page = p
		}
	}

	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			pagination.PageSize = ps
		}
	}

	projects, paginationResp, iterationCounts, err := h.projectUseCase.GetByMemberID(ctx, userID, pagination)
	if err != nil {
		http.Error(w, "Failed to retrieve projects", http.StatusInternalServerError)
		return
	}

	// Add iteration_count to each project in the response
	projectsWithCounts := make([]map[string]interface{}, len(projects))
	for i, project := range projects {
		projectMap := map[string]interface{}{
			"id":               project.ID,
			"name":             project.Name,
			"description":      project.Description,
			"color":            project.Color,
			"prod_range":       project.ProdRange,
			"members":          project.Members,
			"created_at":       project.CreatedAt,
			"updated_at":       project.UpdatedAt,
			"iterations_count": iterationCounts[project.ID],
		}
		projectsWithCounts[i] = projectMap
	}

	response := map[string]interface{}{
		"data":       projectsWithCounts,
		"pagination": paginationResp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
