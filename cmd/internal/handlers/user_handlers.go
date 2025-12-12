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

type UserHandlers struct {
	userUseCase *usecases.UserUseCase
}

func NewUserHandlers(userUseCase *usecases.UserUseCase) *UserHandlers {
	return &UserHandlers{
		userUseCase: userUseCase,
	}
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// GetAllUsers handles GET /users
// @Summary Get all users
// @Description Get a paginated list of all users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20) maximum(100)
// @Success 200 {object} map[string]interface{} "Users with pagination"
// @Failure 500 {string} string "Internal server error"
// @Router /users [get]
func (h *UserHandlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
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

	users, paginationResp, err := h.userUseCase.GetAll(ctx, pagination)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       users,
		"pagination": paginationResp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserByID handles GET /users/{id}
// @Summary Get user by ID
// @Description Get a specific user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID" format(uuid)
// @Success 200 {object} models.User "User details"
// @Failure 400 {string} string "Invalid user ID"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
func (h *UserHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := h.userUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser handles POST /users
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} map[string]interface{} "Created user"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create user"
// @Router /users [post]
func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	newUser := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	ctx := r.Context()
	userID, err := h.userUseCase.Add(ctx, newUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":    userID,
		"name":  req.Name,
		"email": req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateUser handles PUT /users/{id}
// @Summary Update user
// @Description Update an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID" format(uuid)
// @Param user body UpdateUserRequest true "Updated user data"
// @Success 200 {object} models.User "Updated user"
// @Failure 400 {string} string "Invalid user ID or request body"
// @Failure 500 {string} string "Failed to update user"
// @Router /users/{id} [put]
func (h *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	updatedUser := models.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}

	ctx := r.Context()
	err = h.userUseCase.Update(ctx, updatedUser)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser handles DELETE /users/{id}
// @Summary Delete user
// @Description Delete a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID" format(uuid)
// @Success 204 "User deleted successfully"
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Failed to delete user"
// @Router /users/{id} [delete]
func (h *UserHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.userUseCase.Delete(ctx, id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
