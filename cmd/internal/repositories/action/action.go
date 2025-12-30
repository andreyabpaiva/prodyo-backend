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

func (r *Repository) Get(ctx context.Context, indicatorRangeID uuid.UUID) ([]models.Action, error) {
	const query = `
		SELECT a.id, a.indicator_range_id, a.description, a.status, a.start_at, a.end_at, a.assignee_id, a.created_at, a.updated_at,
		       c.id as cause_id, c.metric, c.description as cause_description,
		       c.productivity_level, c.created_at as cause_created_at, c.updated_at as cause_updated_at,
		       u.id as user_id, u.name as user_name, u.email as user_email
		FROM actions a
		INNER JOIN causes c ON a.cause_id = c.id
		LEFT JOIN users u ON a.assignee_id = u.id
		WHERE a.indicator_range_id = $1
		ORDER BY a.created_at ASC
	`
	rows, err := r.db.Query(ctx, query, indicatorRangeID)
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
			&a.IndicatorRangeID,
			&a.Description,
			&a.Status,
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
		c.IndicatorRangeID = a.IndicatorRangeID
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

func (r *Repository) GetByIterationID(ctx context.Context, iterationID uuid.UUID) ([]models.Action, error) {
	const query = `
		SELECT a.id, a.indicator_range_id, a.description, a.status, a.start_at, a.end_at, a.assignee_id, a.created_at, a.updated_at,
		       c.id as cause_id, c.indicator_range_id as cause_indicator_range_id, c.metric, c.description as cause_description,
		       c.productivity_level, c.created_at as cause_created_at, c.updated_at as cause_updated_at,
		       u.id as user_id, u.name as user_name, u.email as user_email
		FROM actions a
		INNER JOIN causes c ON a.cause_id = c.id
		INNER JOIN indicator_ranges ir ON a.indicator_range_id = ir.id
		INNER JOIN iterations i ON ir.project_id = i.project_id
		LEFT JOIN users u ON a.assignee_id = u.id
		WHERE i.id = $1
		ORDER BY a.created_at ASC
	`
	rows, err := r.db.Query(ctx, query, iterationID)
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
		var causeIndicatorRangeID uuid.UUID
		if err := rows.Scan(
			&a.ID,
			&a.IndicatorRangeID,
			&a.Description,
			&a.Status,
			&a.StartAt,
			&a.EndAt,
			&assigneeID,
			&a.CreatedAt,
			&a.UpdatedAt,
			&c.ID,
			&causeIndicatorRangeID,
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
		c.IndicatorRangeID = causeIndicatorRangeID
		a.Cause = c

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

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (models.Action, error) {
	const query = `
		SELECT a.id, a.indicator_range_id, a.description, a.status, a.start_at, a.end_at, a.assignee_id, a.created_at, a.updated_at,
		       c.id as cause_id, c.metric, c.description as cause_description,
		       c.productivity_level, c.created_at as cause_created_at, c.updated_at as cause_updated_at,
		       u.id as user_id, u.name as user_name, u.email as user_email
		FROM actions a
		INNER JOIN causes c ON a.cause_id = c.id
		LEFT JOIN users u ON a.assignee_id = u.id
		WHERE a.id = $1
	`
	var a models.Action
	var c models.Cause
	var assigneeID *uuid.UUID
	var userName, userEmail *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID,
		&a.IndicatorRangeID,
		&a.Description,
		&a.Status,
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
	)
	if err != nil {
		return models.Action{}, ErrNotFound
	}

	c.IndicatorRangeID = a.IndicatorRangeID
	a.Cause = c

	if assigneeID != nil && userName != nil && userEmail != nil {
		a.Assignee = models.User{
			ID:    *assigneeID,
			Name:  *userName,
			Email: *userEmail,
		}
	}

	return a, nil
}

func (r *Repository) Create(ctx context.Context, action models.Action) error {
	const query = `
		INSERT INTO actions (id, indicator_range_id, cause_id, description, status, start_at, end_at, assignee_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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

	status := action.Status
	if status == "" {
		status = models.StatusNotStarted
	}

	_, err := r.db.Exec(ctx, query,
		action.ID,
		action.IndicatorRangeID,
		action.Cause.ID,
		action.Description,
		status,
		startAt,
		endAt,
		assigneeID,
	)
	return err
}

func (r *Repository) Update(ctx context.Context, action models.Action) error {
	const query = `
		UPDATE actions
		SET description = $2, status = $3, start_at = $4, end_at = $5, assignee_id = $6, updated_at = NOW()
		WHERE id = $1
	`
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
		action.Description,
		action.Status,
		startAt,
		endAt,
		assigneeID,
	)
	return err
}
