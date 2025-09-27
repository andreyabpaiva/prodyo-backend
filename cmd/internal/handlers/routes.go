package handlers

import (
	"net/http"
	"prodyo-backend/cmd/internal/usecases"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRoutes configures all the routes for the API
func SetupRoutes(projectUseCase *usecases.ProjectUseCase) *mux.Router {
	router := mux.NewRouter()

	// Initialize handlers
	projectHandlers := NewProjectHandlers(projectUseCase)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Project routes
	api.HandleFunc("/projects", projectHandlers.GetAllProjects).Methods("GET")
	api.HandleFunc("/projects", projectHandlers.CreateProject).Methods("POST")
	api.HandleFunc("/projects/{id}", projectHandlers.GetProjectByID).Methods("GET")
	api.HandleFunc("/projects/{id}", projectHandlers.UpdateProject).Methods("PUT")
	api.HandleFunc("/projects/{id}", projectHandlers.DeleteProject).Methods("DELETE")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// Swagger documentation endpoint
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// CORS middleware
	router.Use(corsMiddleware)

	return router
}

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
