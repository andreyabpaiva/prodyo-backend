package models

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page" form:"page" validate:"min=1"`
	PageSize int `json:"page_size" form:"page_size" validate:"min=1,max=100"`
}

// PaginationResponse represents paginated response metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// GetOffset calculates the offset for pagination
func (p *PaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20 // Default page size
	}
	if p.PageSize > 100 {
		p.PageSize = 100 // Max page size
	}
	return (p.Page - 1) * p.PageSize
}

// NewPaginationResponse creates a new pagination response
func NewPaginationResponse(page, pageSize int, total int64) PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
