package task

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
	ErrNotFound = errors.New("task not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, iterationID uuid.UUID) ([]models.Task, error) {
	const query = `
		SELECT t.id, t.iteration_id, t.name, t.description, t.status, t.timer, t.points, t.parent_task_id,
		       t.created_at, t.updated_at,
		       u.id as assignee_id, u.name as assignee_name, u.email as assignee_email,
		       u.created_at as assignee_created_at, u.updated_at as assignee_updated_at
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
		WHERE t.iteration_id = $1 AND t.parent_task_id IS NULL
		ORDER BY t.created_at ASC
	`
	rows, err := r.db.Query(ctx, query, iterationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		t, err := r.scanTask(rows)
		if err != nil {
			return nil, err
		}
		t.Tasks = []models.Task{}
		t.Improvements = []models.Improv{}
		t.Bugs = []models.Bug{}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (models.Task, error) {
	const query = `
		SELECT t.id, t.iteration_id, t.name, t.description, t.status, t.timer, t.points, t.parent_task_id,
		       t.created_at, t.updated_at,
		       u.id as assignee_id, u.name as assignee_name, u.email as assignee_email,
		       u.created_at as assignee_created_at, u.updated_at as assignee_updated_at
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
		WHERE t.id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	task, err := r.scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Task{}, ErrNotFound
		}
		return models.Task{}, err
	}
	task.Tasks = []models.Task{}
	task.Improvements = []models.Improv{}
	task.Bugs = []models.Bug{}
	return task, nil
}

func (r *Repository) scanTask(row interface {
	Scan(dest ...interface{}) error
}) (models.Task, error) {
	var t models.Task
	var assigneeID *uuid.UUID
	var assigneeName, assigneeEmail *string
	var assigneeCreatedAt, assigneeUpdatedAt *time.Time
	var parentTaskID *uuid.UUID
	var timer *time.Time

	err := row.Scan(
		&t.ID,
		&t.IterationID,
		&t.Name,
		&t.Description,
		&t.Status,
		&timer,
		&t.Points,
		&parentTaskID,
		&t.CreatedAt,
		&t.UpdatedAt,
		&assigneeID,
		&assigneeName,
		&assigneeEmail,
		&assigneeCreatedAt,
		&assigneeUpdatedAt,
	)
	if err != nil {
		return models.Task{}, err
	}

	if assigneeID != nil {
		t.Assignee = models.User{
			ID:        *assigneeID,
			Name:      *assigneeName,
			Email:     *assigneeEmail,
			CreatedAt: *assigneeCreatedAt,
			UpdatedAt: *assigneeUpdatedAt,
		}
	}

	if timer != nil {
		t.Timer = *timer
	}

	return t, nil
}

func (r *Repository) Create(ctx context.Context, task models.Task) error {
	const query = `
		INSERT INTO tasks (id, iteration_id, name, description, assignee_id, status, timer, points, parent_task_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}

	var assigneeID interface{}
	if task.Assignee.ID != uuid.Nil {
		assigneeID = task.Assignee.ID
	} else {
		assigneeID = nil
	}

	var timer interface{}
	if !task.Timer.IsZero() {
		timer = task.Timer
	} else {
		timer = nil
	}

	var parentTaskID interface{}
	// Check if this is a sub-task
	if len(task.Tasks) > 0 {
		// This would be handled differently - sub-tasks are separate records
		parentTaskID = nil
	}

	points := task.Points
	if points == 0 {
		points = 1
	}

	_, err := r.db.Exec(ctx, query,
		task.ID,
		task.IterationID,
		task.Name,
		task.Description,
		assigneeID,
		task.Status,
		timer,
		points,
		parentTaskID,
	)
	return err
}

func (r *Repository) Update(ctx context.Context, task models.Task) error {
	const query = `
		UPDATE tasks
		SET name = $1, description = $2, assignee_id = $3, status = $4, timer = $5, points = $6, updated_at = NOW()
		WHERE id = $7
	`
	var assigneeID interface{}
	if task.Assignee.ID != uuid.Nil {
		assigneeID = task.Assignee.ID
	} else {
		assigneeID = nil
	}

	var timer interface{}
	if !task.Timer.IsZero() {
		timer = task.Timer
	} else {
		timer = nil
	}

	points := task.Points
	if points == 0 {
		points = 1
	}

	cmd, err := r.db.Exec(ctx, query,
		task.Name,
		task.Description,
		assigneeID,
		task.Status,
		timer,
		points,
		task.ID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM tasks WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
