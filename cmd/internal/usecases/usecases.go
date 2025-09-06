package usecases

import (
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories"

	"github.com/google/uuid"
)

type ProjectUseCase struct {
	repos *repositories.Repositories
}

func New(repos *repositories.Repositories) *ProjectUseCase {
	return &ProjectUseCase{
		repos: repos,
	}
}

func (u ProjectUseCase) GetAll() ([]models.Project, error) {
	users, err := u.repos.Project.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u ProjectUseCase) Add(newProject models.Project) uuid.UUID {
	repoReq := models.Project{
		ID:   uuid.New(),
		Name: newProject.Name,
	}

	u.repos.Project.Add(repoReq)

	return repoReq.ID
}
