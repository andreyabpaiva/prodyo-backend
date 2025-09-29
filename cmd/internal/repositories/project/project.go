package project

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Project struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Project {
	return &Project{db: db}
}

func (p Project) GetAll() ([]models.Project, error) {
	const query = `SELECT id, name, email, created_at, updated_at FROM projects ORDER BY created_at DESC`
	rows, err := p.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]models.Project, 0)
	for rows.Next() {
		var pr models.Project
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.Email, &pr.CreatedAt, &pr.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, pr)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return projects, nil
}

func (p Project) GetByID(id uuid.UUID) (models.Project, error) {
	const query = `SELECT id, name, email, created_at, updated_at FROM projects WHERE id = $1`
	var pr models.Project
	err := p.db.QueryRow(context.Background(), query, id).Scan(&pr.ID, &pr.Name, &pr.Email, &pr.CreatedAt, &pr.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Project{}, errors.New("project not found")
		}
		return models.Project{}, err
	}
	return pr, nil
}

func (p *Project) Add(newProject models.Project) error {
	const query = `INSERT INTO projects (id, name, email) VALUES ($1, $2, $3)`
	if newProject.ID == uuid.Nil {
		newProject.ID = uuid.New()
	}
	_, err := p.db.Exec(context.Background(), query, newProject.ID, newProject.Name, newProject.Email)
	return err
}

func (p *Project) Update(project models.Project) error {
	const query = `UPDATE projects SET name = $1, email = $2 WHERE id = $3`
	cmd, err := p.db.Exec(context.Background(), query, project.Name, project.Email, project.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("project not found")
	}
	return nil
}

func (p *Project) Delete(id uuid.UUID) error {
	const query = `DELETE FROM projects WHERE id = $1`
	cmd, err := p.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("project not found")
	}
	return nil
}
