package user

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("user not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, pagination models.PaginationRequest) ([]models.User, models.PaginationResponse, error) {
	countQuery := `SELECT COUNT(*) FROM users`
	var total int64
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	query := `
		SELECT id, name, email, project_id, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, pagination.PageSize, pagination.GetOffset())
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.ProjectID,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, models.PaginationResponse{}, err
		}
		users = append(users, u)
	}

	if rows.Err() != nil {
		return nil, models.PaginationResponse{}, err
	}

	paginationResp := models.NewPaginationResponse(pagination.Page, pagination.PageSize, total)
	return users, paginationResp, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	const query = `
		SELECT id, name, email, project_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.ProjectID,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}
	return u, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	const query = `
		SELECT id, name, email, project_id, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var u models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.ProjectID,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}
	return u, nil
}

func (r *Repository) Add(ctx context.Context, u models.User) error {
	const query = `
		INSERT INTO users (id, name, email)
		VALUES ($1, $2, $3)
	`
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query,
		u.ID,
		u.Name,
		u.Email,
	)
	return err
}

func (r *Repository) Update(ctx context.Context, u models.User) error {
	const query = `
		UPDATE users
		SET name = $1, email = $2, updated_at = NOW()
		WHERE id = $3
	`
	cmd, err := r.db.Exec(ctx, query,
		u.Name,
		u.Email,
		u.ID,
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
	const query = `DELETE FROM users WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
