package handlers

import (
	"net/http"
	"prodyo-backend/cmd/internal/usecases"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(
	projectUseCase *usecases.ProjectUseCase,
	userUseCase *usecases.UserUseCase,
	authUseCase *usecases.AuthUseCase,
	iterationUseCase *usecases.IterationUseCase,
	taskUseCase *usecases.TaskUseCase,
	improvUseCase *usecases.ImprovUseCase,
	bugUseCase *usecases.BugUseCase,
	indicatorUseCase *usecases.IndicatorUseCase,
	indicatorRangeUseCase *usecases.IndicatorRangeUseCase,
	causeUseCase *usecases.CauseUseCase,
	actionUseCase *usecases.ActionUseCase,
) *mux.Router {
	router := mux.NewRouter()

	// Initialize handlers
	projectHandlers := NewProjectHandlers(projectUseCase, indicatorRangeUseCase)
	userHandlers := NewUserHandlers(userUseCase)
	authHandlers := NewAuthHandlers(authUseCase)
	iterationHandlers := NewIterationHandlers(iterationUseCase)
	taskHandlers := NewTaskHandlers(taskUseCase)
	improvHandlers := NewImprovHandlers(improvUseCase)
	bugHandlers := NewBugHandlers(bugUseCase)
	indicatorHandlers := NewIndicatorHandlers(indicatorUseCase, indicatorRangeUseCase, causeUseCase, actionUseCase)

	api := router.PathPrefix("/api/v1").Subrouter()

	// Public routes (no authentication required)
	public := api.PathPrefix("").Subrouter()
	public.HandleFunc("/auth/register", authHandlers.Register).Methods("POST")
	public.HandleFunc("/auth/login", authHandlers.Login).Methods("POST")

	// Protected routes (authentication required)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(AuthMiddleware(authUseCase))

	// Auth routes
	protected.HandleFunc("/auth/logout", authHandlers.Logout).Methods("POST")

	// Project routes
	protected.HandleFunc("/projects", projectHandlers.GetAllProjects).Methods("GET")
	protected.HandleFunc("/projects", projectHandlers.CreateProject).Methods("POST")
	protected.HandleFunc("/projects/member/{userId}", projectHandlers.GetProjectsByMemberID).Methods("GET")
	protected.HandleFunc("/projects/{id}", projectHandlers.GetProjectByID).Methods("GET")
	protected.HandleFunc("/projects/{id}", projectHandlers.UpdateProject).Methods("PUT")
	protected.HandleFunc("/projects/{id}", projectHandlers.DeleteProject).Methods("DELETE")

	// Project indicator ranges routes (project-level)
	protected.HandleFunc("/projects/{project_id}/indicator-ranges", indicatorHandlers.GetRanges).Methods("GET")
	protected.HandleFunc("/projects/{project_id}/indicator-ranges/default", indicatorHandlers.CreateDefaultRanges).Methods("POST")
	protected.HandleFunc("/projects/{project_id}/indicator-ranges/{indicator_type}", indicatorHandlers.GetRangeByIndicatorType).Methods("GET")
	protected.HandleFunc("/projects/{project_id}/indicator-range-ids/{indicator_type}", indicatorHandlers.GetIndicatorRangeIDByProjectIDAndType).Methods("GET")

	// User routes
	protected.HandleFunc("/users", userHandlers.GetAllUsers).Methods("GET")
	protected.HandleFunc("/users", userHandlers.CreateUser).Methods("POST")
	protected.HandleFunc("/users/project/{projectId}", userHandlers.GetUsersByProjectID).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandlers.GetUserByID).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandlers.UpdateUser).Methods("PUT")
	protected.HandleFunc("/users/{id}", userHandlers.DeleteUser).Methods("DELETE")

	// Iteration routes
	protected.HandleFunc("/iterations", iterationHandlers.GetAll).Methods("GET")
	protected.HandleFunc("/iterations", iterationHandlers.Create).Methods("POST")
	protected.HandleFunc("/iterations/{id}", iterationHandlers.GetByID).Methods("GET")
	protected.HandleFunc("/iterations/{id}", iterationHandlers.Delete).Methods("DELETE")
	protected.HandleFunc("/iterations/{id}/analysis", iterationHandlers.GetIterationAnalysis).Methods("GET")
	protected.HandleFunc("/iterations/{iteration_id}/causes-actions", indicatorHandlers.GetCausesAndActionsByIteration).Methods("GET")

	// Task routes
	protected.HandleFunc("/tasks", taskHandlers.GetAll).Methods("GET")
	protected.HandleFunc("/tasks", taskHandlers.Create).Methods("POST")
	protected.HandleFunc("/tasks/{id}", taskHandlers.GetByID).Methods("GET")
	protected.HandleFunc("/tasks/{id}", taskHandlers.Update).Methods("PUT")
	protected.HandleFunc("/tasks/{id}", taskHandlers.Patch).Methods("PATCH")
	protected.HandleFunc("/tasks/{id}", taskHandlers.Delete).Methods("DELETE")

	// Improvement routes
	protected.HandleFunc("/improvements", improvHandlers.GetAll).Methods("GET")
	protected.HandleFunc("/improvements", improvHandlers.Create).Methods("POST")
	protected.HandleFunc("/improvements/{id}", improvHandlers.GetByID).Methods("GET")

	// Bug routes
	protected.HandleFunc("/bugs", bugHandlers.GetAll).Methods("GET")
	protected.HandleFunc("/bugs", bugHandlers.Create).Methods("POST")
	protected.HandleFunc("/bugs/{id}", bugHandlers.GetByID).Methods("GET")

	// Indicator routes
	protected.HandleFunc("/indicators", indicatorHandlers.Get).Methods("GET")
	protected.HandleFunc("/indicators", indicatorHandlers.Create).Methods("POST")
	protected.HandleFunc("/indicators/causes", indicatorHandlers.CreateCause).Methods("POST")
	protected.HandleFunc("/indicators/actions", indicatorHandlers.CreateAction).Methods("POST")
	protected.HandleFunc("/indicators/ranges", indicatorHandlers.SetRange).Methods("POST")
	protected.HandleFunc("/indicators/ranges/{range_id}", indicatorHandlers.DeleteRange).Methods("DELETE")
	protected.HandleFunc("/indicators/{indicator_id}/metrics", indicatorHandlers.UpdateMetricValues).Methods("PUT")
	protected.HandleFunc("/indicators/{indicator_id}/summary", indicatorHandlers.GetMetricSummary).Methods("GET")

	// Health check (public)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// Swagger documentation
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

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

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
