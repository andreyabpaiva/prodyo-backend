package repositories

import (
	"prodyo-backend/cmd/internal/repositories/project"
	"prodyo-backend/cmd/internal/repositories/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Project *project.Repository
	User    *user.Repository
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		Project: project.New(db),
		User:    user.New(db),
	}
}
