package handlers

import (
	"net/http"
	"prodyo-backend/cmd/internal/usecases"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(projectUseCase *usecases.ProjectUseCase, userUseCase *usecases.UserUseCase) *mux.Router {
	router := mux.NewRouter()

	projectHandlers := NewProjectHandlers(projectUseCase)
	userHandlers := NewUserHandlers(userUseCase)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/projects", projectHandlers.GetAllProjects).Methods("GET")
	api.HandleFunc("/projects", projectHandlers.CreateProject).Methods("POST")
	api.HandleFunc("/projects/{id}", projectHandlers.GetProjectByID).Methods("GET")
	api.HandleFunc("/projects/{id}", projectHandlers.UpdateProject).Methods("PUT")
	api.HandleFunc("/projects/{id}", projectHandlers.DeleteProject).Methods("DELETE")

	api.HandleFunc("/users", userHandlers.GetAllUsers).Methods("GET")
	api.HandleFunc("/users", userHandlers.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", userHandlers.GetUserByID).Methods("GET")
	api.HandleFunc("/users/{id}", userHandlers.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandlers.DeleteUser).Methods("DELETE")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.Use(CorsMiddleware)

	return router
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
