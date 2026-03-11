package dto_response

// PaginationMeta standardizes pagination metadata for list responses.
type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	PerPage      int   `json:"per_page"`
}
