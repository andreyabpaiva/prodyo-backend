package handlers

import (
	"encoding/json"
	"net/http"
	"prodyo-backend/cmd/internal/usecases"
)

type AuthHandlers struct {
	authUseCase *usecases.AuthUseCase
}

func NewAuthHandlers(authUseCase *usecases.AuthUseCase) *AuthHandlers {
	return &AuthHandlers{
		authUseCase: authUseCase,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
}

// Register handles POST /auth/register
// @Summary Register a new user
// @Description Create a new user account with email, password, and name
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 409 {string} string "User already exists"
// @Failure 500 {string} string "Failed to register user"
// @Router /auth/register [post]
func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "Email, password, and name are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID, err := h.authUseCase.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":    userID.String(),
		"email": req.Email,
		"name":  req.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles POST /auth/login
// @Summary Login user
// @Description Authenticate user and return session token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Invalid email or password"
// @Router /auth/login [post]
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	session, err := h.authUseCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	user, err := h.authUseCase.ValidateSession(ctx, session.Token)
	if err != nil {
		http.Error(w, "Failed to get user details", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token: session.Token,
	}
	response.User.ID = user.ID.String()
	response.User.Name = user.Name
	response.User.Email = user.Email

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout handles POST /auth/logout
// @Summary Logout user
// @Description Invalidate user session token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 204 "Logout successful"
// @Failure 401 {string} string "Authorization token required"
// @Failure 500 {string} string "Failed to logout"
// @Router /auth/logout [post]
func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	ctx := r.Context()
	if err := h.authUseCase.Logout(ctx, token); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
