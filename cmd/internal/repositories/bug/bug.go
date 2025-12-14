package bug

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
	ErrNotFound = errors.New("bug not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, taskID uuid.UUID) ([]models.Bug, error) {
	const query = `
		SELECT b.id, b.task_id, b.number, b.description, b.points, b.created_at, b.updated_at,
		       u.id as assignee_id, u.name as assignee_name, u.email as assignee_email,
		       u.created_at as assignee_created_at, u.updated_at as assignee_updated_at
		FROM bugs b
		LEFT JOIN users u ON b.assignee_id = u.id
		WHERE b.task_id = $1
		ORDER BY b.number ASC
	`
	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bugs []models.Bug
	for rows.Next() {
		var bg models.Bug
		var assigneeID *uuid.UUID
		var assigneeName, assigneeEmail *string
		var assigneeCreatedAt, assigneeUpdatedAt *time.Time

		if err := rows.Scan(
			&bg.ID,
			&bg.TaskID,
			&bg.Number,
			&bg.Description,
			&bg.Points,
			&bg.CreatedAt,
			&bg.UpdatedAt,
			&assigneeID,
			&assigneeName,
			&assigneeEmail,
			&assigneeCreatedAt,
			&assigneeUpdatedAt,
		); err != nil {
			return nil, err
		}

		if assigneeID != nil {
			bg.Assignee = models.User{
				ID:        *assigneeID,
				Name:      *assigneeName,
				Email:     *assigneeEmail,
				CreatedAt: *assigneeCreatedAt,
				UpdatedAt: *assigneeUpdatedAt,
			}
		}

		bugs = append(bugs, bg)
	}

	return bugs, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (models.Bug, error) {
	const query = `
		SELECT b.id, b.task_id, b.number, b.description, b.points, b.created_at, b.updated_at,
		       u.id as assignee_id, u.name as assignee_name, u.email as assignee_email,
		       u.created_at as assignee_created_at, u.updated_at as assignee_updated_at
		FROM bugs b
		LEFT JOIN users u ON b.assignee_id = u.id
		WHERE b.id = $1
	`
	var bg models.Bug
	var assigneeID *uuid.UUID
	var assigneeName, assigneeEmail *string
	var assigneeCreatedAt, assigneeUpdatedAt *time.Time

	err := r.db.QueryRow(ctx, query, id).Scan(
		&bg.ID,
		&bg.TaskID,
		&bg.Number,
		&bg.Description,
		&bg.Points,
		&bg.CreatedAt,
		&bg.UpdatedAt,
		&assigneeID,
		&assigneeName,
		&assigneeEmail,
		&assigneeCreatedAt,
		&assigneeUpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Bug{}, ErrNotFound
		}
		return models.Bug{}, err
	}

	if assigneeID != nil {
		bg.Assignee = models.User{
			ID:        *assigneeID,
			Name:      *assigneeName,
			Email:     *assigneeEmail,
			CreatedAt: *assigneeCreatedAt,
			UpdatedAt: *assigneeUpdatedAt,
		}
	}

	return bg, nil
}

func (r *Repository) Create(ctx context.Context, bug models.Bug) error {
	const query = `
		INSERT INTO bugs (id, task_id, assignee_id, number, description, points)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if bug.ID == uuid.Nil {
		bug.ID = uuid.New()
	}

	var assigneeID interface{}
	if bug.Assignee.ID != uuid.Nil {
		assigneeID = bug.Assignee.ID
	} else {
		assigneeID = nil
	}

	points := bug.Points
	if points == 0 {
		points = 1
	}

	_, err := r.db.Exec(ctx, query,
		bug.ID,
		bug.TaskID,
		assigneeID,
		bug.Number,
		bug.Description,
		points,
	)
	return err
}
