package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/task"

	"github.com/google/uuid"
)

type TaskUseCase struct {
	repo *task.Repository
}

func NewTaskUseCase(repo *task.Repository) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (u *TaskUseCase) GetAll(ctx context.Context, iterationID uuid.UUID) ([]models.Task, error) {
	return u.repo.GetAll(ctx, iterationID)
}

func (u *TaskUseCase) GetByID(ctx context.Context, id uuid.UUID) (models.Task, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *TaskUseCase) Create(ctx context.Context, newTask models.Task) (uuid.UUID, error) {
	if newTask.ID == uuid.Nil {
		newTask.ID = uuid.New()
	}

	if err := u.repo.Create(ctx, newTask); err != nil {
		return uuid.Nil, err
	}

	return newTask.ID, nil
}

func (u *TaskUseCase) Update(ctx context.Context, task models.Task) error {
	return u.repo.Update(ctx, task)
}

func (u *TaskUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

