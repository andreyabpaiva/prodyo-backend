package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/project"

	"github.com/google/uuid"
)

type ProjectUseCase struct {
	repo *project.Repository
}

// Construtor
func NewProjectUseCase(repo *project.Repository) *ProjectUseCase {
	return &ProjectUseCase{repo: repo}
}

func (u *ProjectUseCase) GetAll(ctx context.Context, pagination models.PaginationRequest) ([]models.Project, models.PaginationResponse, error) {
	return u.repo.GetAll(ctx, pagination)
}

func (u *ProjectUseCase) GetByID(ctx context.Context, id uuid.UUID) (models.Project, int64, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *ProjectUseCase) Add(ctx context.Context, newProject models.Project) (uuid.UUID, error) {
	if newProject.ID == uuid.Nil {
		newProject.ID = uuid.New()
	}

	if err := u.repo.Add(ctx, newProject); err != nil {
		return uuid.Nil, err
	}

	return newProject.ID, nil
}

func (u *ProjectUseCase) Update(ctx context.Context, project models.Project) error {
	return u.repo.Update(ctx, project)
}

func (u *ProjectUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func (u *ProjectUseCase) GetByMemberID(ctx context.Context, userID uuid.UUID, pagination models.PaginationRequest) ([]models.Project, models.PaginationResponse, map[uuid.UUID]int64, error) {
	return u.repo.GetByMemberID(ctx, userID, pagination)
}