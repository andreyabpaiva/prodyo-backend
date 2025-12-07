package repositories

import (
	"prodyo-backend/cmd/internal/repositories/action"
	"prodyo-backend/cmd/internal/repositories/bug"
	"prodyo-backend/cmd/internal/repositories/cause"
	"prodyo-backend/cmd/internal/repositories/indicator"
	"prodyo-backend/cmd/internal/repositories/iteration"
	"prodyo-backend/cmd/internal/repositories/improv"
	"prodyo-backend/cmd/internal/repositories/project"
	"prodyo-backend/cmd/internal/repositories/session"
	"prodyo-backend/cmd/internal/repositories/task"
	"prodyo-backend/cmd/internal/repositories/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Project    *project.Repository
	User       *user.Repository
	Session    *session.Repository
	Iteration  *iteration.Repository
	Task       *task.Repository
	Improv     *improv.Repository
	Bug        *bug.Repository
	Indicator  *indicator.Repository
	Cause      *cause.Repository
	Action     *action.Repository
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		Project:   project.New(db),
		User:      user.New(db),
		Session:   session.New(db),
		Iteration: iteration.New(db),
		Task:      task.New(db),
		Improv:    improv.New(db),
		Bug:       bug.New(db),
		Indicator: indicator.New(db),
		Cause:     cause.New(db),
		Action:    action.New(db),
	}
}
