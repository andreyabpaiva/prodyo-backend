package repositories

import (
	"prodyo-backend/cmd/internal/repositories/project"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Project *project.Repository
	// Adicione outros repositórios aqui
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		Project: project.New(db),
	}
}
