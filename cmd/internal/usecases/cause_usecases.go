package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/cause"

	"github.com/google/uuid"
)

type CauseUseCase struct {
	repo *cause.Repository
}

func NewCauseUseCase(repo *cause.Repository) *CauseUseCase {
	return &CauseUseCase{repo: repo}
}

func (u *CauseUseCase) Get(ctx context.Context, indicatorID uuid.UUID) ([]models.Cause, error) {
	return u.repo.Get(ctx, indicatorID)
}

func (u *CauseUseCase) Create(ctx context.Context, cause models.Cause) (uuid.UUID, error) {
	if cause.ID == uuid.Nil {
		cause.ID = uuid.New()
	}

	if err := u.repo.Create(ctx, cause); err != nil {
		return uuid.Nil, err
	}

	return cause.ID, nil
}

