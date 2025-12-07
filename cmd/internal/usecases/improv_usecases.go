package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/improv"

	"github.com/google/uuid"
)

type ImprovUseCase struct {
	repo *improv.Repository
}

func NewImprovUseCase(repo *improv.Repository) *ImprovUseCase {
	return &ImprovUseCase{repo: repo}
}

func (u *ImprovUseCase) GetAll(ctx context.Context, taskID uuid.UUID) ([]models.Improv, error) {
	return u.repo.GetAll(ctx, taskID)
}

func (u *ImprovUseCase) GetByID(ctx context.Context, id uuid.UUID) (models.Improv, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *ImprovUseCase) Create(ctx context.Context, improv models.Improv) (uuid.UUID, error) {
	if improv.ID == uuid.Nil {
		improv.ID = uuid.New()
	}

	if err := u.repo.Create(ctx, improv); err != nil {
		return uuid.Nil, err
	}

	return improv.ID, nil
}

