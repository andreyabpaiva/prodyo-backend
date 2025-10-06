package project

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("project not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context) ([]models.Project, error) {
	const query = `
		SELECT id, name, description, color, prod_range, email, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var pr models.Project
		if err := rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.Description,
			&pr.Color,
			&pr.ProdRange,
			&pr.Email,
			&pr.CreatedAt,
			&pr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, pr)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return projects, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (models.Project, error) {
	const query = `
		SELECT id, name, description, color, prod_range, email, created_at, updated_at
		FROM projects
		WHERE id = $1
	`
	var pr models.Project
	err := r.db.QueryRow(ctx, query, id).Scan(
		&pr.ID,
		&pr.Name,
		&pr.Description,
		&pr.Color,
		&pr.ProdRange,
		&pr.Email,
		&pr.CreatedAt,
		&pr.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Project{}, ErrNotFound
		}
		return models.Project{}, err
	}
	return pr, nil
}

func (r *Repository) Add(ctx context.Context, pr models.Project) error {
	const query = `
		INSERT INTO projects (id, name, description, color, prod_range, email)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if pr.ID == uuid.Nil {
		pr.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query,
		pr.ID,
		pr.Name,
		pr.Description,
		pr.Color,
		pr.ProdRange,
		pr.Email,
	)
	return err
}

func (r *Repository) Update(ctx context.Context, pr models.Project) error {
	const query = `
		UPDATE projects
		SET name = $1, description = $2, color = $3, prod_range = $4, email = $5, updated_at = NOW()
		WHERE id = $6
	`
	cmd, err := r.db.Exec(ctx, query,
		pr.Name,
		pr.Description,
		pr.Color,
		pr.ProdRange,
		pr.Email,
		pr.ID,
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
	const query = `DELETE FROM projects WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
