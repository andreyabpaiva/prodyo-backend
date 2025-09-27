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
	projects, err := u.repos.Project.GetAll()
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (u ProjectUseCase) GetByID(id uuid.UUID) (models.Project, error) {
	project, err := u.repos.Project.GetByID(id)
	if err != nil {
		return models.Project{}, err
	}

	return project, nil
}

func (u ProjectUseCase) Add(newProject models.Project) uuid.UUID {
	repoReq := models.Project{
		ID:    uuid.New(),
		Name:  newProject.Name,
		Email: newProject.Email,
	}

	u.repos.Project.Add(repoReq)

	return repoReq.ID
}

func (u ProjectUseCase) Update(project models.Project) error {
	err := u.repos.Project.Update(project)
	if err != nil {
		return err
	}

	return nil
}

func (u ProjectUseCase) Delete(id uuid.UUID) error {
	err := u.repos.Project.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
