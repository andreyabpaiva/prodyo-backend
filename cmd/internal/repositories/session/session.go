package session

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
	ErrNotFound = errors.New("session not found")
	ErrExpired  = errors.New("session expired")
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, session models.Session) error {
	const query = `
		INSERT INTO sessions (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}

	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.Token,
		session.ExpiresAt,
	)
	return err
}

func (r *Repository) GetByToken(ctx context.Context, token string) (models.Session, error) {
	const query = `
		SELECT id, user_id, token, expires_at, created_at
		FROM sessions
		WHERE token = $1
	`
	var s models.Session
	err := r.db.QueryRow(ctx, query, token).Scan(
		&s.ID,
		&s.UserID,
		&s.Token,
		&s.ExpiresAt,
		&s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Session{}, ErrNotFound
		}
		return models.Session{}, err
	}

	// Check if session is expired
	if time.Now().After(s.ExpiresAt) {
		return models.Session{}, ErrExpired
	}

	return s, nil
}

func (r *Repository) Delete(ctx context.Context, token string) error {
	const query = `DELETE FROM sessions WHERE token = $1`
	cmd, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	const query = `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *Repository) CleanExpired(ctx context.Context) error {
	const query = `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := r.db.Exec(ctx, query)
	return err
}

