package indicator_range

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("indicator range not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (models.IndicatorRange, error) {
	const query = `
		SELECT id, project_id, indicator_type,
			ok_min, ok_max, alert_min, alert_max, critical_min, critical_max,
			created_at, updated_at
		FROM indicator_ranges
		WHERE id = $1
	`

	var ir models.IndicatorRange
	var okMin, okMax, alertMin, alertMax, criticalMin, criticalMax float64

	err := r.db.QueryRow(ctx, query, id).Scan(
		&ir.ID,
		&ir.ProjectID,
		&ir.IndicatorType,
		&okMin, &okMax,
		&alertMin, &alertMax,
		&criticalMin, &criticalMax,
		&ir.CreatedAt,
		&ir.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.IndicatorRange{}, ErrNotFound
		}
		return models.IndicatorRange{}, err
	}

	ir.Range = models.ProductivityRange{
		Ok:       models.RangeValues{Min: okMin, Max: okMax},
		Alert:    models.RangeValues{Min: alertMin, Max: alertMax},
		Critical: models.RangeValues{Min: criticalMin, Max: criticalMax},
	}

	return ir, nil
}

// GetByProjectID returns all indicator ranges for a project
func (r *Repository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.IndicatorRange, error) {
	const query = `
		SELECT id, project_id, indicator_type,
			ok_min, ok_max, alert_min, alert_max, critical_min, critical_max,
			created_at, updated_at
		FROM indicator_ranges
		WHERE project_id = $1
		ORDER BY indicator_type
	`

	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ranges []models.IndicatorRange
	for rows.Next() {
		var ir models.IndicatorRange
		var okMin, okMax, alertMin, alertMax, criticalMin, criticalMax float64

		err := rows.Scan(
			&ir.ID,
			&ir.ProjectID,
			&ir.IndicatorType,
			&okMin, &okMax,
			&alertMin, &alertMax,
			&criticalMin, &criticalMax,
			&ir.CreatedAt,
			&ir.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		ir.Range = models.ProductivityRange{
			Ok:       models.RangeValues{Min: okMin, Max: okMax},
			Alert:    models.RangeValues{Min: alertMin, Max: alertMax},
			Critical: models.RangeValues{Min: criticalMin, Max: criticalMax},
		}

		ranges = append(ranges, ir)
	}

	if ranges == nil {
		ranges = []models.IndicatorRange{}
	}

	return ranges, rows.Err()
}

// GetByIndicatorType returns the range for a specific indicator type in a project
func (r *Repository) GetByIndicatorType(ctx context.Context, projectID uuid.UUID, indicatorType models.IndicatorEnum) (models.IndicatorRange, error) {
	const query = `
		SELECT id, project_id, indicator_type,
			ok_min, ok_max, alert_min, alert_max, critical_min, critical_max,
			created_at, updated_at
		FROM indicator_ranges
		WHERE project_id = $1 AND indicator_type = $2
	`

	var ir models.IndicatorRange
	var okMin, okMax, alertMin, alertMax, criticalMin, criticalMax float64

	err := r.db.QueryRow(ctx, query, projectID, indicatorType).Scan(
		&ir.ID,
		&ir.ProjectID,
		&ir.IndicatorType,
		&okMin, &okMax,
		&alertMin, &alertMax,
		&criticalMin, &criticalMax,
		&ir.CreatedAt,
		&ir.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.IndicatorRange{}, ErrNotFound
		}
		return models.IndicatorRange{}, err
	}

	ir.Range = models.ProductivityRange{
		Ok:       models.RangeValues{Min: okMin, Max: okMax},
		Alert:    models.RangeValues{Min: alertMin, Max: alertMax},
		Critical: models.RangeValues{Min: criticalMin, Max: criticalMax},
	}

	return ir, nil
}

// Create creates or updates an indicator range (upsert)
func (r *Repository) Create(ctx context.Context, ir models.IndicatorRange) error {
	const query = `
		INSERT INTO indicator_ranges (id, project_id, indicator_type, ok_min, ok_max, alert_min, alert_max, critical_min, critical_max)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (project_id, indicator_type) DO UPDATE SET
			ok_min = EXCLUDED.ok_min,
			ok_max = EXCLUDED.ok_max,
			alert_min = EXCLUDED.alert_min,
			alert_max = EXCLUDED.alert_max,
			critical_min = EXCLUDED.critical_min,
			critical_max = EXCLUDED.critical_max,
			updated_at = NOW()
	`

	if ir.ID == uuid.Nil {
		ir.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query,
		ir.ID,
		ir.ProjectID,
		ir.IndicatorType,
		ir.Range.Ok.Min, ir.Range.Ok.Max,
		ir.Range.Alert.Min, ir.Range.Alert.Max,
		ir.Range.Critical.Min, ir.Range.Critical.Max,
	)
	return err
}

// CreateDefaultRanges creates default indicator ranges for a new project
func (r *Repository) CreateDefaultRanges(ctx context.Context, projectID uuid.UUID) error {
	// Default ranges for each indicator type
	defaults := []models.IndicatorRange{
		{
			ID:            uuid.New(),
			ProjectID:     projectID,
			IndicatorType: models.IndicatorSpeedPerIteration,
			Range: models.ProductivityRange{
				Ok:       models.RangeValues{Min: 10, Max: 100},
				Alert:    models.RangeValues{Min: 5, Max: 10},
				Critical: models.RangeValues{Min: 0, Max: 5},
			},
		},
		{
			ID:            uuid.New(),
			ProjectID:     projectID,
			IndicatorType: models.IndicatorReworkPerIteration,
			Range: models.ProductivityRange{
				Ok:       models.RangeValues{Min: 0, Max: 0.1},
				Alert:    models.RangeValues{Min: 0.1, Max: 0.3},
				Critical: models.RangeValues{Min: 0.3, Max: 1},
			},
		},
		{
			ID:            uuid.New(),
			ProjectID:     projectID,
			IndicatorType: models.IndicatorInstabilityIndex,
			Range: models.ProductivityRange{
				Ok:       models.RangeValues{Min: 0, Max: 0.15},
				Alert:    models.RangeValues{Min: 0.15, Max: 0.4},
				Critical: models.RangeValues{Min: 0.4, Max: 1},
			},
		},
	}

	for _, ir := range defaults {
		if err := r.Create(ctx, ir); err != nil {
			return err
		}
	}

	return nil
}

// Update updates an existing indicator range
func (r *Repository) Update(ctx context.Context, ir models.IndicatorRange) error {
	const query = `
		UPDATE indicator_ranges
		SET ok_min = $2, ok_max = $3,
			alert_min = $4, alert_max = $5,
			critical_min = $6, critical_max = $7,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		ir.ID,
		ir.Range.Ok.Min, ir.Range.Ok.Max,
		ir.Range.Alert.Min, ir.Range.Alert.Max,
		ir.Range.Critical.Min, ir.Range.Critical.Max,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes an indicator range
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM indicator_ranges WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteByProjectID removes all indicator ranges for a project
func (r *Repository) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	const query = `DELETE FROM indicator_ranges WHERE project_id = $1`
	_, err := r.db.Exec(ctx, query, projectID)
	return err
}
