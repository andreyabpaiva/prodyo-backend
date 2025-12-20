package indicator

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("indicator not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, iterationID uuid.UUID) (models.Indicator, error) {
	const query = `
		SELECT id, iteration_id, 
			COALESCE(velocity_value, 0), COALESCE(rework_value, 0), COALESCE(instability_value, 0),
			created_at, updated_at
		FROM indicators
		WHERE iteration_id = $1
	`
	var ind models.Indicator
	err := r.db.QueryRow(ctx, query, iterationID).Scan(
		&ind.ID,
		&ind.IterationID,
		&ind.SpeedValue,
		&ind.ReworkValue,
		&ind.InstabilityValue,
		&ind.CreatedAt,
		&ind.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Indicator{}, ErrNotFound
		}
		return models.Indicator{}, err
	}
	ind.Causes = []models.Cause{}
	ind.Actions = []models.Action{}
	return ind, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (models.Indicator, error) {
	const query = `
		SELECT id, iteration_id, 
			COALESCE(velocity_value, 0), COALESCE(rework_value, 0), COALESCE(instability_value, 0),
			created_at, updated_at
		FROM indicators
		WHERE id = $1
	`
	var ind models.Indicator
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ind.ID,
		&ind.IterationID,
		&ind.SpeedValue,
		&ind.ReworkValue,
		&ind.InstabilityValue,
		&ind.CreatedAt,
		&ind.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Indicator{}, ErrNotFound
		}
		return models.Indicator{}, err
	}
	ind.Causes = []models.Cause{}
	ind.Actions = []models.Action{}
	return ind, nil
}

func (r *Repository) Create(ctx context.Context, indicator models.Indicator) error {
	const query = `
		INSERT INTO indicators (id, iteration_id, velocity_value, rework_value, instability_value)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (iteration_id) DO NOTHING
	`
	if indicator.ID == uuid.Nil {
		indicator.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query,
		indicator.ID,
		indicator.IterationID,
		indicator.SpeedValue,
		indicator.ReworkValue,
		indicator.InstabilityValue,
	)
	return err
}

func (r *Repository) UpdateMetricValues(ctx context.Context, indicatorID uuid.UUID, speed, rework, instability float64) error {
	const query = `
		UPDATE indicators
		SET velocity_value = $2, rework_value = $3, instability_value = $4, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, indicatorID, speed, rework, instability)
	return err
}

// GetProjectIDByIterationID retrieves the project_id for an iteration
// This is useful for getting the project-level indicator ranges
func (r *Repository) GetProjectIDByIterationID(ctx context.Context, iterationID uuid.UUID) (uuid.UUID, error) {
	const query = `
		SELECT project_id FROM iterations WHERE id = $1
	`
	var projectID uuid.UUID
	err := r.db.QueryRow(ctx, query, iterationID).Scan(&projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrNotFound
		}
		return uuid.Nil, err
	}
	return projectID, nil
}
