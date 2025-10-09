package project

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"
	"time"

	"github.com/google/uuid"
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

func (r *Repository) GetAll(ctx context.Context, pagination models.PaginationRequest) ([]models.Project, models.PaginationResponse, error) {
	// First, get total count
	countQuery := `SELECT COUNT(DISTINCT p.id) FROM projects p`
	var total int64
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	// Then get paginated results
	query := `
		SELECT 
			p.id, p.name, p.description, p.color, p.prod_range, p.created_at, p.updated_at,
			u.id as member_id, u.name as member_name, u.email as member_email, 
			u.created_at as member_created_at, u.updated_at as member_updated_at
		FROM projects p
		LEFT JOIN project_members pm ON p.id = pm.project_id
		LEFT JOIN users u ON pm.user_id = u.id
		ORDER BY p.created_at DESC, u.name ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, pagination.PageSize, pagination.GetOffset())
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}
	defer rows.Close()

	projectMap := make(map[uuid.UUID]*models.Project)

	for rows.Next() {
		var pr models.Project
		var memberID *uuid.UUID
		var memberName, memberEmail *string
		var memberCreatedAt, memberUpdatedAt *time.Time

		if err := rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.Description,
			&pr.Color,
			&pr.ProdRange,
			&pr.CreatedAt,
			&pr.UpdatedAt,
			&memberID,
			&memberName,
			&memberEmail,
			&memberCreatedAt,
			&memberUpdatedAt,
		); err != nil {
			return nil, models.PaginationResponse{}, err
		}

		// Get or create project in map
		if existingProject, exists := projectMap[pr.ID]; exists {
			pr = *existingProject
		} else {
			pr.Members = []models.User{}
			projectMap[pr.ID] = &pr
		}

		// Add member if exists
		if memberID != nil {
			member := models.User{
				ID:        *memberID,
				Name:      *memberName,
				Email:     *memberEmail,
				CreatedAt: *memberCreatedAt,
				UpdatedAt: *memberUpdatedAt,
			}
			pr.Members = append(pr.Members, member)
		}
	}

	if rows.Err() != nil {
		return nil, models.PaginationResponse{}, rows.Err()
	}

	// Convert map to slice
	var projects []models.Project
	for _, project := range projectMap {
		projects = append(projects, *project)
	}

	paginationResp := models.NewPaginationResponse(pagination.Page, pagination.PageSize, total)
	return projects, paginationResp, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (models.Project, error) {
	const query = `
		SELECT 
			p.id, p.name, p.description, p.color, p.prod_range, p.created_at, p.updated_at,
			u.id as member_id, u.name as member_name, u.email as member_email, 
			u.created_at as member_created_at, u.updated_at as member_updated_at
		FROM projects p
		LEFT JOIN project_members pm ON p.id = pm.project_id
		LEFT JOIN users u ON pm.user_id = u.id
		WHERE p.id = $1
		ORDER BY u.name ASC
	`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return models.Project{}, err
	}
	defer rows.Close()

	var pr models.Project
	var found bool

	for rows.Next() {
		if !found {
			var memberID *uuid.UUID
			var memberName, memberEmail *string
			var memberCreatedAt, memberUpdatedAt *time.Time

			err := rows.Scan(
				&pr.ID,
				&pr.Name,
				&pr.Description,
				&pr.Color,
				&pr.ProdRange,
				&pr.CreatedAt,
				&pr.UpdatedAt,
				&memberID,
				&memberName,
				&memberEmail,
				&memberCreatedAt,
				&memberUpdatedAt,
			)
			if err != nil {
				return models.Project{}, err
			}

			pr.Members = []models.User{}
			found = true

			// Add member if exists
			if memberID != nil {
				member := models.User{
					ID:        *memberID,
					Name:      *memberName,
					Email:     *memberEmail,
					CreatedAt: *memberCreatedAt,
					UpdatedAt: *memberUpdatedAt,
				}
				pr.Members = append(pr.Members, member)
			}
		} else {
			// Additional members for the same project
			var memberID *uuid.UUID
			var memberName, memberEmail *string
			var memberCreatedAt, memberUpdatedAt *time.Time

			err := rows.Scan(
				&pr.ID,
				&pr.Name,
				&pr.Description,
				&pr.Color,
				&pr.ProdRange,
				&pr.CreatedAt,
				&pr.UpdatedAt,
				&memberID,
				&memberName,
				&memberEmail,
				&memberCreatedAt,
				&memberUpdatedAt,
			)
			if err != nil {
				return models.Project{}, err
			}

			if memberID != nil {
				member := models.User{
					ID:        *memberID,
					Name:      *memberName,
					Email:     *memberEmail,
					CreatedAt: *memberCreatedAt,
					UpdatedAt: *memberUpdatedAt,
				}
				pr.Members = append(pr.Members, member)
			}
		}
	}

	if !found {
		return models.Project{}, ErrNotFound
	}

	return pr, nil
}

func (r *Repository) Add(ctx context.Context, pr models.Project) error {
	const query = `
		INSERT INTO projects (id, name, description, color, prod_range)
		VALUES ($1, $2, $3, $4, $5)
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
	)
	if err != nil {
		return err
	}

	// Add members to the project
	return r.addProjectMembers(ctx, pr.ID, pr.Members)
}

func (r *Repository) Update(ctx context.Context, pr models.Project) error {
	const query = `
		UPDATE projects
		SET name = $1, description = $2, color = $3, prod_range = $4, updated_at = NOW()
		WHERE id = $5
	`
	cmd, err := r.db.Exec(ctx, query,
		pr.Name,
		pr.Description,
		pr.Color,
		pr.ProdRange,
		pr.ID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}

	// Update members for the project
	if err := r.updateProjectMembers(ctx, pr.ID, pr.Members); err != nil {
		return err
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

// Helper methods for managing project members
func (r *Repository) getProjectMembers(ctx context.Context, projectID uuid.UUID) ([]models.User, error) {
	const query = `
		SELECT u.id, u.name, u.email, u.created_at, u.updated_at
		FROM users u
		INNER JOIN project_members pm ON u.id = pm.user_id
		WHERE pm.project_id = $1
		ORDER BY u.name
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, u)
	}

	return members, rows.Err()
}

func (r *Repository) addProjectMembers(ctx context.Context, projectID uuid.UUID, members []models.User) error {
	if len(members) == 0 {
		return nil
	}

	const query = `
		INSERT INTO project_members (project_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (project_id, user_id) DO NOTHING
	`

	for _, member := range members {
		_, err := r.db.Exec(ctx, query, projectID, member.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) updateProjectMembers(ctx context.Context, projectID uuid.UUID, members []models.User) error {
	// First, remove all existing members
	const deleteQuery = `DELETE FROM project_members WHERE project_id = $1`
	_, err := r.db.Exec(ctx, deleteQuery, projectID)
	if err != nil {
		return err
	}

	// Then add the new members
	return r.addProjectMembers(ctx, projectID, members)
}
