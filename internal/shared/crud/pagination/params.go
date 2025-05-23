package pagination

// FilterMap represents a key-value filter map
type FilterMap map[string]string

// QueryParams represents parameters for pagination and filtering
type QueryParams struct {
	Page   int       `form:"page"`
	Limit  int       `form:"limit"`
	Sort   string    `form:"sort"`
	Order  string    `form:"order"`
	Filter FilterMap `form:"-"`
	Search string    `form:"q"`
}

// PaginationParams extracts pagination parameters with defaults
func (q *QueryParams) PaginationParams() (page, limit int, sort, order string) {
	page = q.Page
	if page < 1 {
		page = 1
	}

	limit = q.Limit
	if limit <= 0 {
		limit = 20 // Default limit
	}

	sort = q.Sort
	if sort == "" {
		sort = "created_at" // Default sort field
	}

	order = q.Order
	if order == "" {
		order = "desc" // Default order
	}

	return page, limit, sort, order
}
