package action

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("action not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, indicatorID uuid.UUID) ([]models.Action, error) {
	const query = `
		SELECT a.id, a.indicator_id, a.description, a.created_at, a.updated_at,
		       c.id as cause_id, c.metric, c.description as cause_description,
		       c.productivity_level, c.created_at as cause_created_at, c.updated_at as cause_updated_at
		FROM actions a
		INNER JOIN causes c ON a.cause_id = c.id
		WHERE a.indicator_id = $1
		ORDER BY a.created_at ASC
	`
	rows, err := r.db.Query(ctx, query, indicatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []models.Action
	for rows.Next() {
		var a models.Action
		var c models.Cause
		if err := rows.Scan(
			&a.ID,
			&a.IndicatorID,
			&a.Description,
			&a.CreatedAt,
			&a.UpdatedAt,
			&c.ID,
			&c.Metric,
			&c.Description,
			&c.ProductivityLevel,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		c.IndicatorID = a.IndicatorID
		a.Cause = c
		actions = append(actions, a)
	}

	return actions, rows.Err()
}

func (r *Repository) Create(ctx context.Context, action models.Action) error {
	const query = `
		INSERT INTO actions (id, indicator_id, cause_id, description)
		VALUES ($1, $2, $3, $4)
	`
	if action.ID == uuid.Nil {
		action.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query,
		action.ID,
		action.IndicatorID,
		action.Cause.ID,
		action.Description,
	)
	return err
}

