package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/bug"

	"github.com/google/uuid"
)

type BugUseCase struct {
	repo *bug.Repository
}

func NewBugUseCase(repo *bug.Repository) *BugUseCase {
	return &BugUseCase{repo: repo}
}

func (u *BugUseCase) GetAll(ctx context.Context, taskID uuid.UUID) ([]models.Bug, error) {
	return u.repo.GetAll(ctx, taskID)
}

func (u *BugUseCase) GetByID(ctx context.Context, id uuid.UUID) (models.Bug, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *BugUseCase) Create(ctx context.Context, bug models.Bug) (uuid.UUID, error) {
	if bug.ID == uuid.Nil {
		bug.ID = uuid.New()
	}

	if err := u.repo.Create(ctx, bug); err != nil {
		return uuid.Nil, err
	}

	return bug.ID, nil
}

