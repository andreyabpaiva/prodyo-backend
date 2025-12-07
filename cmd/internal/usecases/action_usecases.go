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

func (u *ActionUseCase) Create(ctx context.Context, action models.Action) (uuid.UUID, error) {
	if action.ID == uuid.Nil {
		action.ID = uuid.New()
	}

	if err := u.repo.Create(ctx, action); err != nil {
		return uuid.Nil, err
	}

	return action.ID, nil
}

