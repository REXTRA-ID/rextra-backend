package utils

import "math"

// PaginationParams holds incoming pagination request values.
type PaginationParams struct {
	Page  int
	Limit int
}

// PaginationMeta standardizes pagination metadata in responses.
type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	PerPage      int   `json:"per_page"`
}

// NewPaginationParams normalizes page and limit into safe bounds (page>=1, limit 1..100).
func NewPaginationParams(page, limit int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// GetOffset returns the offset for use in SQL limit/offset queries.
func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// CalculatePaginationMeta builds meta info from page, limit, and total rows.
func CalculatePaginationMeta(page, limit int, total int64) PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return PaginationMeta{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalRecords: total,
		PerPage:      limit,
	}
}
