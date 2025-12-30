package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/action"

	"github.com/google/uuid"
)

type ActionUseCase struct {
	repo *action.Repository
}

func NewActionUseCase(repo *action.Repository) *ActionUseCase {
	return &ActionUseCase{repo: repo}
}

func (u *ActionUseCase) Get(ctx context.Context, indicatorID uuid.UUID) ([]models.Action, error) {
	return u.repo.Get(ctx, indicatorID)
}

func (u *ActionUseCase) GetByIterationID(ctx context.Context, iterationID uuid.UUID) ([]models.Action, error) {
	return u.repo.GetByIterationID(ctx, iterationID)
}

func (u *ActionUseCase) GetByID(ctx context.Context, id uuid.UUID) (models.Action, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *ActionUseCase) Create(ctx context.Context, action models.Action) (uuid.UUID, error) {
	if action.ID == uuid.Nil {
		action.ID = uuid.New()
	}

	if err := u.repo.Create(ctx, action); err != nil {
		return uuid.Nil, err
	}

	return action.ID, nil
}

func (u *ActionUseCase) Update(ctx context.Context, action models.Action) error {
	return u.repo.Update(ctx, action)
}

