package dto

type PaginatedResultDTO[T any] struct {
	Result     []T   `json:"result"`      // The slice of items for the current page
	TotalCount int64 `json:"total_count"` // Total number of items across all pages
	Page       int   `json:"page"`        // Current page number
	Limit      int   `json:"limit"`       // Number of items per page
	TotalPages int   `json:"total_pages"` // Total number of pages
	HasNext    bool  `json:"has_next"`    // Whether there are more pages after this one
	HasPrev    bool  `json:"has_prev"`    // Whether there are pages before this one
}
