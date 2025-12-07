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
		SELECT id, iteration_id, created_at, updated_at
		FROM indicators
		WHERE iteration_id = $1
	`
	var ind models.Indicator
	err := r.db.QueryRow(ctx, query, iterationID).Scan(
		&ind.ID,
		&ind.IterationID,
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
		INSERT INTO indicators (id, iteration_id)
		VALUES ($1, $2)
		ON CONFLICT (iteration_id) DO NOTHING
	`
	if indicator.ID == uuid.Nil {
		indicator.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query,
		indicator.ID,
		indicator.IterationID,
	)
	return err
}

