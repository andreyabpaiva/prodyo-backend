package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/indicator_range"
	"prodyo-backend/cmd/internal/repositories/iteration"
	"prodyo-backend/cmd/internal/repositories/task"
	"prodyo-backend/cmd/internal/services"

	"github.com/google/uuid"
)

type IterationUseCase struct {
	repo               *iteration.Repository
	taskRepo           *task.Repository
	indicatorRangeRepo *indicator_range.Repository
}

func NewIterationUseCase(repo *iteration.Repository, taskRepo *task.Repository, indicatorRangeRepo *indicator_range.Repository) *IterationUseCase {
	return &IterationUseCase{
		repo:               repo,
		taskRepo:           taskRepo,
		indicatorRangeRepo: indicatorRangeRepo,
	}
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

func (u *IterationUseCase) GetIterationAnalysis(ctx context.Context, iterationID uuid.UUID) (models.IterationAnalysisResponse, error) {
	iteration, err := u.repo.GetByID(ctx, iterationID)
	if err != nil {
		return models.IterationAnalysisResponse{}, err
	}

	tasks, err := u.taskRepo.GetAll(ctx, iterationID)
	if err != nil {
		return models.IterationAnalysisResponse{}, err
	}

	ranges, err := u.indicatorRangeRepo.GetByProjectID(ctx, iteration.ProjectID)
	if err != nil {
		return models.IterationAnalysisResponse{}, err
	}

	calculator := services.NewIndicatorCalculator(tasks, ranges)
	analysis := calculator.CalculateIterationAnalysis(iterationID)

	return analysis, nil
}
