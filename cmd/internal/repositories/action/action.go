package action

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"
	"time"

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
		SELECT a.id, a.indicator_id, a.description, a.start_at, a.end_at, a.assignee_id, a.created_at, a.updated_at,
		       c.id as cause_id, c.metric, c.description as cause_description,
		       c.productivity_level, c.created_at as cause_created_at, c.updated_at as cause_updated_at,
		       u.id as user_id, u.name as user_name, u.email as user_email
		FROM actions a
		INNER JOIN causes c ON a.cause_id = c.id
		LEFT JOIN users u ON a.assignee_id = u.id
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
		var assigneeID *uuid.UUID
		var userName, userEmail *string
		if err := rows.Scan(
			&a.ID,
			&a.IndicatorID,
			&a.Description,
			&a.StartAt,
			&a.EndAt,
			&assigneeID,
			&a.CreatedAt,
			&a.UpdatedAt,
			&c.ID,
			&c.Metric,
			&c.Description,
			&c.ProductivityLevel,
			&c.CreatedAt,
			&c.UpdatedAt,
			&assigneeID,
			&userName,
			&userEmail,
		); err != nil {
			return nil, err
		}
		c.IndicatorID = a.IndicatorID
		a.Cause = c

		// Populate assignee if present
		if assigneeID != nil && userName != nil && userEmail != nil {
			a.Assignee = models.User{
				ID:    *assigneeID,
				Name:  *userName,
				Email: *userEmail,
			}
		}

		actions = append(actions, a)
	}

	return actions, rows.Err()
}

func (r *Repository) Create(ctx context.Context, action models.Action) error {
	const query = `
		INSERT INTO actions (id, indicator_id, cause_id, description, start_at, end_at, assignee_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	if action.ID == uuid.Nil {
		action.ID = uuid.New()
	}

	var assigneeID *uuid.UUID
	if action.Assignee.ID != uuid.Nil {
		assigneeID = &action.Assignee.ID
	}

	var startAt, endAt *time.Time
	if !action.StartAt.IsZero() {
		startAt = &action.StartAt
	}
	if !action.EndAt.IsZero() {
		endAt = &action.EndAt
	}

	_, err := r.db.Exec(ctx, query,
		action.ID,
		action.IndicatorID,
		action.Cause.ID,
		action.Description,
		startAt,
		endAt,
		assigneeID,
	)
	return err
}

