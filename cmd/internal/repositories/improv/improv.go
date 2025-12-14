package improv

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("improvement not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, taskID uuid.UUID) ([]models.Improv, error) {
	const query = `
		SELECT i.id, i.task_id, i.number, i.description, i.points, i.created_at, i.updated_at,
		       u.id as assignee_id, u.name as assignee_name, u.email as assignee_email,
		       u.created_at as assignee_created_at, u.updated_at as assignee_updated_at
		FROM improvements i
		LEFT JOIN users u ON i.assignee_id = u.id
		WHERE i.task_id = $1
		ORDER BY i.number ASC
	`
	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var improvements []models.Improv
	for rows.Next() {
		var imp models.Improv
		var assigneeID *uuid.UUID
		var assigneeName, assigneeEmail *string
		var assigneeCreatedAt, assigneeUpdatedAt *time.Time

		if err := rows.Scan(
			&imp.ID,
			&imp.TaskID,
			&imp.Number,
			&imp.Description,
			&imp.Points,
			&imp.CreatedAt,
			&imp.UpdatedAt,
			&assigneeID,
			&assigneeName,
			&assigneeEmail,
			&assigneeCreatedAt,
			&assigneeUpdatedAt,
		); err != nil {
			return nil, err
		}

		if assigneeID != nil {
			imp.Assignee = models.User{
				ID:        *assigneeID,
				Name:      *assigneeName,
				Email:     *assigneeEmail,
				CreatedAt: *assigneeCreatedAt,
				UpdatedAt: *assigneeUpdatedAt,
			}
		}

		improvements = append(improvements, imp)
	}

	return improvements, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (models.Improv, error) {
	const query = `
		SELECT i.id, i.task_id, i.number, i.description, i.points, i.created_at, i.updated_at,
		       u.id as assignee_id, u.name as assignee_name, u.email as assignee_email,
		       u.created_at as assignee_created_at, u.updated_at as assignee_updated_at
		FROM improvements i
		LEFT JOIN users u ON i.assignee_id = u.id
		WHERE i.id = $1
	`
	var imp models.Improv
	var assigneeID *uuid.UUID
	var assigneeName, assigneeEmail *string
	var assigneeCreatedAt, assigneeUpdatedAt *time.Time

	err := r.db.QueryRow(ctx, query, id).Scan(
		&imp.ID,
		&imp.TaskID,
		&imp.Number,
		&imp.Description,
		&imp.Points,
		&imp.CreatedAt,
		&imp.UpdatedAt,
		&assigneeID,
		&assigneeName,
		&assigneeEmail,
		&assigneeCreatedAt,
		&assigneeUpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Improv{}, ErrNotFound
		}
		return models.Improv{}, err
	}

	if assigneeID != nil {
		imp.Assignee = models.User{
			ID:        *assigneeID,
			Name:      *assigneeName,
			Email:     *assigneeEmail,
			CreatedAt: *assigneeCreatedAt,
			UpdatedAt: *assigneeUpdatedAt,
		}
	}

	return imp, nil
}

func (r *Repository) Create(ctx context.Context, improv models.Improv) error {
	const query = `
		INSERT INTO improvements (id, task_id, assignee_id, number, description, points)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if improv.ID == uuid.Nil {
		improv.ID = uuid.New()
	}

	var assigneeID interface{}
	if improv.Assignee.ID != uuid.Nil {
		assigneeID = improv.Assignee.ID
	} else {
		assigneeID = nil
	}

	points := improv.Points
	if points == 0 {
		points = 1
	}

	_, err := r.db.Exec(ctx, query,
		improv.ID,
		improv.TaskID,
		assigneeID,
		improv.Number,
		improv.Description,
		points,
	)
	return err
}
