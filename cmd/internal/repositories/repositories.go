package repositories

import (
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/project"

	"github.com/google/uuid"
)

type Repositories struct {
	Project interface {
		GetAll() ([]models.Project, error)
		GetByID(id uuid.UUID) (models.Project, error)
		Add(newProject models.Project) error
		Update(project models.Project) error
		Delete(id uuid.UUID) error
	}
}

func New() *Repositories {
	return &Repositories{
		Project: project.New(),
	}
}
