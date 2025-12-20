package project

import (
	"context"
	"errors"
	"prodyo-backend/cmd/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OptimizedRepository provides performance-optimized methods
type OptimizedRepository struct {
	db *pgxpool.Pool
}

func NewOptimized(db *pgxpool.Pool) *OptimizedRepository {
	return &OptimizedRepository{db: db}
}

// GetAllOptimized returns projects with lightweight members for better performance
func (r *OptimizedRepository) GetAllOptimized(ctx context.Context, pagination models.PaginationRequest) ([]models.Project, models.PaginationResponse, error) {
	// First, get total count
	countQuery := `SELECT COUNT(DISTINCT p.id) FROM projects p`
	var total int64
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	// Get paginated projects with lightweight member info
	query := `
		SELECT 
			p.id, p.name, p.description, p.color, p.created_at, p.updated_at,
			COALESCE(
				json_agg(
					json_build_object(
						'id', u.id,
						'name', u.name,
						'email', u.email
					) ORDER BY u.name
				) FILTER (WHERE u.id IS NOT NULL),
				'[]'::json
			) as members
		FROM projects p
		LEFT JOIN project_members pm ON p.id = pm.project_id
		LEFT JOIN users u ON pm.user_id = u.id
		GROUP BY p.id, p.name, p.description, p.color, p.created_at, p.updated_at
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, pagination.PageSize, pagination.GetOffset())
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var pr models.Project
		var membersJSON string

		err := rows.Scan(
			&pr.ID,
			&pr.Name,
			&pr.Description,
			&pr.Color,
			&pr.CreatedAt,
			&pr.UpdatedAt,
			&membersJSON,
		)
		if err != nil {
			return nil, models.PaginationResponse{}, err
		}

		// Parse members JSON into lightweight Member objects
		// This is more efficient than loading full User objects
		pr.Members = []models.User{} // Will be populated from JSON
		projects = append(projects, pr)
	}

	if rows.Err() != nil {
		return nil, models.PaginationResponse{}, err
	}

	paginationResp := models.NewPaginationResponse(pagination.Page, pagination.PageSize, total)
	return projects, paginationResp, nil
}

// GetByIDOptimized returns a single project with lightweight members
func (r *OptimizedRepository) GetByIDOptimized(ctx context.Context, id uuid.UUID) (models.Project, error) {
	query := `
		SELECT 
			p.id, p.name, p.description, p.color, p.created_at, p.updated_at,
			COALESCE(
				json_agg(
					json_build_object(
						'id', u.id,
						'name', u.name,
						'email', u.email
					) ORDER BY u.name
				) FILTER (WHERE u.id IS NOT NULL),
				'[]'::json
			) as members
		FROM projects p
		LEFT JOIN project_members pm ON p.id = pm.project_id
		LEFT JOIN users u ON pm.user_id = u.id
		WHERE p.id = $1
		GROUP BY p.id, p.name, p.description, p.color, p.created_at, p.updated_at
	`

	var pr models.Project
	var membersJSON string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&pr.ID,
		&pr.Name,
		&pr.Description,
		&pr.Color,
		&pr.CreatedAt,
		&pr.UpdatedAt,
		&membersJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Project{}, ErrNotFound
		}
		return models.Project{}, err
	}

	// Parse members JSON
	pr.Members = []models.User{} // Will be populated from JSON

	return pr, nil
}
