package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/iteration"

	"github.com/google/uuid"
)

type IterationUseCase struct {
	repo *iteration.Repository
}

func NewIterationUseCase(repo *iteration.Repository) *IterationUseCase {
	return &IterationUseCase{repo: repo}
}

func (u *IterationUseCase) GetAll(ctx context.Context, projectID uuid.UUID) ([]models.Iteration, error) {
	return u.repo.GetAll(ctx, projectID)
}

func (u *IterationUseCase) GetByID(ctx context.Context, id uuid.UUID) (models.Iteration, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *IterationUseCase) Create(ctx context.Context, iteration models.Iteration) (uuid.UUID, error) {
	if iteration.ID == uuid.Nil {
		iteration.ID = uuid.New()
	}

	if err := u.repo.Create(ctx, iteration); err != nil {
		return uuid.Nil, err
	}

	return iteration.ID, nil
}

func (u *IterationUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

