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

// @securityDefinitions.basic BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
package main

import (
	"log"
	"net/http"
	"prodyo-backend/cmd/internal/config"
	"prodyo-backend/cmd/internal/handlers"
	"prodyo-backend/cmd/internal/repositories"
	"prodyo-backend/cmd/internal/usecases"
	_ "prodyo-backend/docs"
)

func main() {
	cfg := config.Load()

	db := repositories.NewDB(cfg.DSN())
	defer db.Close()

	repos := repositories.New()

	projectUseCase := usecases.New(repos)

	router := handlers.SetupRoutes(projectUseCase)

	// log.Println("📋 Available endpoints:")
	// log.Println("  GET    /api/v1/projects     - Get all projects")
	// log.Println("  POST   /api/v1/projects     - Create a new project")
	// log.Println("  GET    /api/v1/projects/{id} - Get project by ID")
	// log.Println("  PUT    /api/v1/projects/{id} - Update project by ID")
	// log.Println("  DELETE /api/v1/projects/{id} - Delete project by ID")
	// log.Println("  GET    /health              - Health check")
	// log.Println("  GET    /swagger/index.html  - Swagger UI documentation")

	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
