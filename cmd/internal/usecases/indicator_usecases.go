package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/indicator"
	"prodyo-backend/cmd/internal/repositories/indicator_range"

	"github.com/google/uuid"
)

type IndicatorUseCase struct {
	repo      *indicator.Repository
	rangeRepo *indicator_range.Repository
}

func NewIndicatorUseCase(repo *indicator.Repository, rangeRepo *indicator_range.Repository) *IndicatorUseCase {
	return &IndicatorUseCase{
		repo:      repo,
		rangeRepo: rangeRepo,
	}
}

func (u *IndicatorUseCase) Get(ctx context.Context, iterationID uuid.UUID) (models.Indicator, error) {
	ind, err := u.repo.Get(ctx, iterationID)
	if err != nil {
		return models.Indicator{}, err
	}

	// Get project ID to fetch ranges
	projectID, err := u.repo.GetProjectIDByIterationID(ctx, iterationID)
	if err == nil {
		// Load project-level ranges and calculate productivity levels
		ranges, err := u.rangeRepo.GetByProjectID(ctx, projectID)
		if err == nil && len(ranges) > 0 {
			ind.CalculateProductivityLevels(ranges)
		}
	}

	return ind, nil
}

func (u *IndicatorUseCase) GetByID(ctx context.Context, id uuid.UUID) (models.Indicator, error) {
	ind, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return models.Indicator{}, err
	}

	// Get project ID to fetch ranges
	projectID, err := u.repo.GetProjectIDByIterationID(ctx, ind.IterationID)
	if err == nil {
		// Load project-level ranges and calculate productivity levels
		ranges, err := u.rangeRepo.GetByProjectID(ctx, projectID)
		if err == nil && len(ranges) > 0 {
			ind.CalculateProductivityLevels(ranges)
		}
	}

	return ind, nil
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

// UpdateMetricValues updates the calculated metric values for an indicator
func (u *IndicatorUseCase) UpdateMetricValues(ctx context.Context, indicatorID uuid.UUID, speed, rework, instability float64) error {
	return u.repo.UpdateMetricValues(ctx, indicatorID, speed, rework, instability)
}

// CalculateAndUpdateMetrics calculates the productivity metrics based on iteration data
// speed = completed tasks / time (in days)
// rework = bugs / tasks
// instability = improvements / tasks
func (u *IndicatorUseCase) CalculateAndUpdateMetrics(ctx context.Context, indicatorID uuid.UUID, totalTasks, completedTasks, bugs, improvements int, durationDays float64) error {
	var speed, rework, instability float64

	// Speed: completed tasks per day
	if durationDays > 0 {
		speed = float64(completedTasks) / durationDays
	}

	// Rework index: bugs per task
	if totalTasks > 0 {
		rework = float64(bugs) / float64(totalTasks)
	}

	// Instability index: improvements per task
	if totalTasks > 0 {
		instability = float64(improvements) / float64(totalTasks)
	}

	return u.repo.UpdateMetricValues(ctx, indicatorID, speed, rework, instability)
}

// GetIndicatorWithLevels returns the indicator with productivity levels calculated
func (u *IndicatorUseCase) GetIndicatorWithLevels(ctx context.Context, iterationID uuid.UUID) (models.Indicator, error) {
	return u.Get(ctx, iterationID)
}

// GetProjectIDByIterationID returns the project ID for an iteration
func (u *IndicatorUseCase) GetProjectIDByIterationID(ctx context.Context, iterationID uuid.UUID) (uuid.UUID, error) {
	return u.repo.GetProjectIDByIterationID(ctx, iterationID)
}
