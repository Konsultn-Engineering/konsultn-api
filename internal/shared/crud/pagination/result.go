package pagination

// PaginatedResult provides a standardized structure for paginated queries
type PaginatedResult[T any] struct {
	Result     []T   `json:"result"`      // The slice of items for the current page
	TotalCount int64 `json:"total_count"` // Total number of items across all pages
	Page       int   `json:"page"`        // Current page number
	Limit      int   `json:"limit"`       // Number of items per page
	TotalPages int   `json:"total_pages"` // Total number of pages
	HasNext    bool  `json:"has_next"`    // Whether there are more pages after this one
	HasPrev    bool  `json:"has_prev"`    // Whether there are pages before this one
}

// NewPaginatedResult creates a new paginated result from the given parameters
func NewPaginatedResult[T any](result []T, totalCount int64, page, limit int) *PaginatedResult[T] {
	// Calculate total pages
	totalPages := int(totalCount) / limit
	if int(totalCount)%limit > 0 {
		totalPages++
	}

	return &PaginatedResult[T]{
		Result:     result,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
