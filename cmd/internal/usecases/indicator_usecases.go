package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/indicator"

	"github.com/google/uuid"
)

type IndicatorUseCase struct {
	repo *indicator.Repository
}

func NewIndicatorUseCase(repo *indicator.Repository) *IndicatorUseCase {
	return &IndicatorUseCase{repo: repo}
}

func (u *IndicatorUseCase) Get(ctx context.Context, iterationID uuid.UUID) (models.Indicator, error) {
	return u.repo.Get(ctx, iterationID)
}

func (u *IndicatorUseCase) Create(ctx context.Context, indicator models.Indicator) (uuid.UUID, error) {
	if indicator.ID == uuid.Nil {
		indicator.ID = uuid.New()
	}

	if err := u.repo.Create(ctx, indicator); err != nil {
		return uuid.Nil, err
	}

	return indicator.ID, nil
}

