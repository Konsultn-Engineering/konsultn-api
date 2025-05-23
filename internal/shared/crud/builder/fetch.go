package builder

import (
	"fmt"
	"gorm.io/gorm"
	"konsultn-api/internal/shared/crud/pagination"
	"konsultn-api/internal/shared/crud/types"
	"strings"
)

// fetch is an internal helper method that retrieves query results
// It can return results either as model instances or as maps
// Parameters:
//   - asMap: Whether to return results as maps (true) or model instances (false)
//
// Returns:
//   - interface{}: Either []T or []map[string]interface{} depending on asMap
//   - error: Any error that occurred during the fetch operation
func (qb *QueryBuilder[T]) fetch(asMap bool) (interface{}, error) {
	if asMap {
		var results []map[string]interface{}
		err := qb.build().Find(&results).Error
		return results, err
	}

	var models []T
	err := qb.build().Find(&models).Error
	return models, err
}

// First retrieves the first record that matches the query conditions
// This is typically used when you expect a single result or want just the first match
// Returns:
//   - *T: Pointer to the model instance, or nil if no match
//   - error: Any error that occurred, including gorm.ErrRecordNotFound if no record found
func (qb *QueryBuilder[T]) First() (*T, error) {
	var model T
	err := qb.build().First(&model).Error
	return &model, err
}

// FirstAsMap retrieves the first record that matches the query conditions as a map
// This is useful when you don't need a full model instance or want dynamic field access
// Returns:
//   - map[string]interface{}: The record as a key-value map, keys being column names
//   - error: Any error that occurred, including gorm.ErrRecordNotFound if no record found
func (qb *QueryBuilder[T]) FirstAsMap() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := qb.build().First(&result).Error
	return result, err
}

// All retrieves all records that match the query conditions as model instances
// This returns the complete result set without pagination
// Returns:
//   - []T: Slice of model instances
//   - error: Any error that occurred during the fetch operation
func (qb *QueryBuilder[T]) All() ([]T, error) {
	result, err := qb.fetch(false)
	if err != nil {
		return nil, err
	}
	return result.([]T), nil
}

// AllAsMaps retrieves all records that match the query conditions as maps
// This is useful for dynamic field access or when you don't need model instances
// Returns:
//   - []map[string]interface{}: Slice of records as key-value maps
//   - error: Any error that occurred during the fetch operation
func (qb *QueryBuilder[T]) AllAsMaps() ([]map[string]interface{}, error) {
	result, err := qb.fetch(true)
	if err != nil {
		return nil, err
	}
	return result.([]map[string]interface{}), nil
}

// Into scans the query results into a custom destination structure
// This is useful when you want to map to a different struct than the model type
// Parameters:
//   - dest: Pointer to the destination structure or slice to scan into
//
// Returns:
//   - error: Any error that occurred during scanning
func (qb *QueryBuilder[T]) Into(dest interface{}) error {
	return qb.build().Scan(dest).Error
}

// WithPageParams configures the query builder with pagination parameters
// This sets up sorting, ordering, page number, and items per page
// Parameters:
//   - params: An object implementing the pagination.QueryParams interface
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WithPageParams(params pagination.QueryParams) types.QueryBuilder[T] {
	page, limit, sortStr, orderStr := params.PaginationParams()

	qb.page = page
	qb.limit = limit

	// Parse sort fields and order directions
	sortFields := strings.Split(sortStr, ",")
	orderFields := strings.Split(orderStr, ",")

	var orders []string
	for i := 0; i < len(sortFields); i++ {
		field := strings.TrimSpace(sortFields[i])

		// If no dot notation is present, prefix with main table name
		if !strings.Contains(field, ".") && !qb.knownAliases[field] {
			tableName := qb.baseTable
			field = tableName + "." + field
		}

		orderDir := "asc" // default order direction
		if i < len(orderFields) && (strings.ToLower(orderFields[i]) == "asc" || strings.ToLower(orderFields[i]) == "desc") {
			orderDir = orderFields[i]
		}

		orders = append(orders, fmt.Sprintf("%s %s", field, orderDir))
	}

	qb.orderClause = strings.Join(orders, ", ")

	return qb
}

// preparePagination sets up pagination and returns the prepared query and total count
func (qb *QueryBuilder[T]) preparePagination() (*gorm.DB, int64, error) {
	var totalCount int64

	// Build the base query
	db := qb.build()

	// Clone DB query for counting
	countDB := db.Session(&gorm.Session{})
	if err := countDB.Count(&totalCount).Order("").Error; err != nil {
		return nil, 0, err
	}

	db = db.Order(qb.orderClause)
	// Apply pagination
	offset := (qb.page - 1) * qb.limit
	db = db.Offset(offset).Limit(qb.limit)

	return db, totalCount, nil
}

// Paginate executes the query with pagination and returns a paginated result
// This method applies the configured page and limit settings, and includes total count
// Returns:
//   - *pagination.PaginatedResult[T]: Structure containing results, total count, and pagination info
//   - error: Any error that occurred during query execution
func (qb *QueryBuilder[T]) Paginate() (*pagination.PaginatedResult[T], error) {
	db, totalCount, err := qb.preparePagination()
	if err != nil {
		return nil, err
	}

	// Execute query with proper type
	var results []T
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	// Create paginated result with the correctly typed slice
	return pagination.NewPaginatedResult(results, totalCount, qb.page, qb.limit), nil
}

// PaginateMap executes the query with pagination and populates a map
func (qb *QueryBuilder[T]) PaginateMap() (pagination.PaginatedResult[map[string]interface{}], error) {
	db, totalCount, err := qb.preparePagination()
	if err != nil {
		return pagination.PaginatedResult[map[string]interface{}]{}, err
	}

	// Create a slice to hold the map results
	var results []map[string]interface{}

	// Execute query directly into slice of maps
	if err := db.Find(&results).Error; err != nil {
		return pagination.PaginatedResult[map[string]interface{}]{}, err
	}

	return *pagination.NewPaginatedResult[map[string]interface{}](results, totalCount, qb.page, qb.limit), nil
}
