package usecases

import (
	"context"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/repositories/indicator_range"

	"github.com/google/uuid"
)

type IndicatorRangeUseCase struct {
	repo *indicator_range.Repository
}

func NewIndicatorRangeUseCase(repo *indicator_range.Repository) *IndicatorRangeUseCase {
	return &IndicatorRangeUseCase{repo: repo}
}

// GetByProjectID returns all indicator ranges for a project
func (u *IndicatorRangeUseCase) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.IndicatorRange, error) {
	return u.repo.GetByProjectID(ctx, projectID)
}

// GetByIndicatorType returns the range for a specific indicator type in a project
func (u *IndicatorRangeUseCase) GetByIndicatorType(ctx context.Context, projectID uuid.UUID, indicatorType models.IndicatorEnum) (models.IndicatorRange, error) {
	return u.repo.GetByIndicatorType(ctx, projectID, indicatorType)
}

// SetRange creates or updates an indicator range for a project
func (u *IndicatorRangeUseCase) SetRange(ctx context.Context, ir models.IndicatorRange) (uuid.UUID, error) {
	if ir.ID == uuid.Nil {
		ir.ID = uuid.New()
	}

	if err := u.repo.Create(ctx, ir); err != nil {
		return uuid.Nil, err
	}

	return ir.ID, nil
}

// CreateDefaultRanges creates default indicator ranges for a new project
func (u *IndicatorRangeUseCase) CreateDefaultRanges(ctx context.Context, projectID uuid.UUID) error {
	return u.repo.CreateDefaultRanges(ctx, projectID)
}

// Update updates an existing indicator range
func (u *IndicatorRangeUseCase) Update(ctx context.Context, ir models.IndicatorRange) error {
	return u.repo.Update(ctx, ir)
}

// Delete removes an indicator range
func (u *IndicatorRangeUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

// DeleteByProjectID removes all indicator ranges for a project
func (u *IndicatorRangeUseCase) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	return u.repo.DeleteByProjectID(ctx, projectID)
}

