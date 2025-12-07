package cause

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("cause not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, indicatorID uuid.UUID) ([]models.Cause, error) {
	const query = `
		SELECT id, indicator_id, metric, description, productivity_level, created_at, updated_at
		FROM causes
		WHERE indicator_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, indicatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var causes []models.Cause
	for rows.Next() {
		var c models.Cause
		if err := rows.Scan(
			&c.ID,
			&c.IndicatorID,
			&c.Metric,
			&c.Description,
			&c.ProductivityLevel,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		causes = append(causes, c)
	}

	return causes, rows.Err()
}

func (r *Repository) Create(ctx context.Context, cause models.Cause) error {
	const query = `
		INSERT INTO causes (id, indicator_id, metric, description, productivity_level)
		VALUES ($1, $2, $3, $4, $5)
	`
	if cause.ID == uuid.Nil {
		cause.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query,
		cause.ID,
		cause.IndicatorID,
		cause.Metric,
		cause.Description,
		cause.ProductivityLevel,
	)
	return err
}

