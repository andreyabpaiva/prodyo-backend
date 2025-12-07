// @title Prodyo Backend API
// @version 1.0
// @description A REST API for managing projects
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
package main

import (
	"log"
	"net/http"
	"prodyo-backend/cmd/internal/config"
	"prodyo-backend/cmd/internal/handlers"
	"prodyo-backend/cmd/internal/migrations"
	"prodyo-backend/cmd/internal/repositories"
	"prodyo-backend/cmd/internal/usecases"
	_ "prodyo-backend/docs"
)

func main() {
	cfg := config.Load()

	log.Println("Running database migrations...")
	if err := migrations.RunMigrations(cfg.DSN(), "cmd/migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	db := repositories.NewDB(cfg.DSN())
	defer db.Close()

	repos := repositories.New(db)

	// Initialize use cases
	projectUseCase := usecases.NewProjectUseCase(repos.Project)
	userUseCase := usecases.NewUserUseCase(repos.User)
	authUseCase := usecases.NewAuthUseCase(repos.User, repos.Session)
	iterationUseCase := usecases.NewIterationUseCase(repos.Iteration)
	taskUseCase := usecases.NewTaskUseCase(repos.Task)
	improvUseCase := usecases.NewImprovUseCase(repos.Improv)
	bugUseCase := usecases.NewBugUseCase(repos.Bug)
	indicatorUseCase := usecases.NewIndicatorUseCase(repos.Indicator)
	causeUseCase := usecases.NewCauseUseCase(repos.Cause)
	actionUseCase := usecases.NewActionUseCase(repos.Action)

	router := handlers.SetupRoutes(
		projectUseCase,
		userUseCase,
		authUseCase,
		iterationUseCase,
		taskUseCase,
		improvUseCase,
		bugUseCase,
		indicatorUseCase,
		causeUseCase,
		actionUseCase,
	)

	handler := handlers.CorsMiddleware(router)

	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
