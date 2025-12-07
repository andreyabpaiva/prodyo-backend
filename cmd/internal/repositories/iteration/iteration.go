package iteration

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("iteration not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, projectID uuid.UUID) ([]models.Iteration, error) {
	const query = `
		SELECT id, project_id, number, description, start_at, end_at, created_at, updated_at
		FROM iterations
		WHERE project_id = $1
		ORDER BY number ASC
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var iterations []models.Iteration
	for rows.Next() {
		var it models.Iteration
		if err := rows.Scan(
			&it.ID,
			&it.ProjectID,
			&it.Number,
			&it.Description,
			&it.StartAt,
			&it.EndAt,
			&it.CreatedAt,
			&it.UpdatedAt,
		); err != nil {
			return nil, err
		}
		it.Tasks = []models.Task{}
		iterations = append(iterations, it)
	}

	return iterations, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (models.Iteration, error) {
	const query = `
		SELECT id, project_id, number, description, start_at, end_at, created_at, updated_at
		FROM iterations
		WHERE id = $1
	`
	var it models.Iteration
	err := r.db.QueryRow(ctx, query, id).Scan(
		&it.ID,
		&it.ProjectID,
		&it.Number,
		&it.Description,
		&it.StartAt,
		&it.EndAt,
		&it.CreatedAt,
		&it.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Iteration{}, ErrNotFound
		}
		return models.Iteration{}, err
	}
	it.Tasks = []models.Task{}
	return it, nil
}

func (r *Repository) Create(ctx context.Context, iteration models.Iteration) error {
	const query = `
		INSERT INTO iterations (id, project_id, number, description, start_at, end_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if iteration.ID == uuid.Nil {
		iteration.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query,
		iteration.ID,
		iteration.ProjectID,
		iteration.Number,
		iteration.Description,
		iteration.StartAt,
		iteration.EndAt,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM iterations WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

