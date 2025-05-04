package crud

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type iRepository[T any, ID idType] interface {
	iCRUDRepository[T, ID]
}

// Repository provides a comprehensive set of CRUD operations over type T with primary key ID.
type Repository[T any, ID idType] struct {
	db           *gorm.DB
	selectFields []string
}

// NewRepository creates and returns a new instance of Repository[T].
func NewRepository[T any, ID idType](db *gorm.DB) *Repository[T, ID] {
	return &Repository[T, ID]{db: db}
}

// Select sets the fields to be selected in the query.
func (r *Repository[T, ID]) Select(fields []string) *Repository[T, ID] {
	newRepo := r.clone()
	if len(fields) > 0 {
		newRepo.db = newRepo.db.Select(fields)
	}
	return newRepo
}

// Count returns total records in the table.
func (r *Repository[T, ID]) Count() (int64, error) {
	var count int64
	err := r.db.Count(&count).Error
	return count, err
}

func (r *Repository[T, ID]) FindWithJoins(
	joins []JoinClause,
	where []WhereClause,
	preloads []string,
) ([]T, error) {
	var result []T
	db := r.db.Model(new(T))

	// Add joins
	for _, join := range joins {
		db = db.Joins(fmt.Sprintf("%s %s ON %s", join.JoinType, join.Table, join.On))
	}

	// Add where conditions
	for _, whereClause := range where {
		db = db.Where(whereClause.Query, whereClause.Args...)
	}

	// Add preloads
	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	err := db.Find(&result).Error
	return result, err
}

func (r *Repository[T, ID]) Query(query AdvancedQuery) ([]T, error) {
	var models []T
	db := r.db.Model(new(T))

	// Apply joins
	for _, join := range query.Joins {
		db = db.Joins(fmt.Sprintf("%s %s ON %s", join.JoinType, join.Table, join.On))
	}

	// Apply where clauses
	for _, where := range query.Wheres {
		db = db.Where(where.Query, where.Args...)
	}

	// Apply filters from QueryParams
	fmt.Println(query.Filter)
	for key, value := range query.Filter {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Apply preloads
	for _, preload := range query.Preload {
		db = db.Preload(preload)
	}

	// Pagination
	offset := (query.Page - 1) * query.Limit
	if offset < 0 {
		offset = 0
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}

	// Sorting
	sortField := "created_at"
	if query.Sort != "" {
		sortField = query.Sort
	}
	order := "desc"
	if strings.ToLower(query.Order) == "asc" {
		order = "asc"
	}

	err := db.Order(fmt.Sprintf("%s %s", sortField, order)).
		Offset(offset).
		Limit(query.Limit).
		Find(&models).Error

	return models, err
}

// FindAll retrieves all records for the given model type.
func (r *Repository[T, ID]) FindAll() ([]T, error) {
	models := make([]T, 0)
	err := r.db.Find(&models).Error
	return models, err
}

func (r *Repository[T, ID]) List(params QueryParams) ([]T, error) {
	var models []T
	db := r.db

	// Apply filters
	for key, value := range params.Filter {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Pagination
	offset := (params.Page - 1) * params.Limit
	if offset < 0 {
		offset = 0
	}

	if params.Limit <= 0 {
		params.Limit = 20
	}

	sortField := "created_at"
	if params.Sort != "" {
		sortField = params.Sort
	}

	order := "desc"
	if strings.ToLower(params.Order) == "asc" {
		order = "asc"
	}

	err := db.Order(fmt.Sprintf("%s %s", sortField, order)).
		Offset(offset).
		Limit(params.Limit).
		Find(&models).Error

	return models, err
}

// FindWhere retrieves records from the database where all provided field-value pairs match exactly.
// Useful for filtering by multiple columns with AND conditions.
// Example: FindWhere(map[string]interface{}{"team_id": "abc", "role": "owner"})
func (r *Repository[T, ID]) FindWhere(filters map[string]interface{}) ([]T, error) {
	var models []T
	err := r.db.Where(filters).Find(&models).Error
	return models, err
}

// FindWhereExpr retrieves records using a custom WHERE clause and arguments.
// Allows more complex queries (e.g., OR conditions, IN clauses, LIKE, etc.).
// Example: FindWhereExpr("status = ? OR priority = ?", "active", "high")
func (r *Repository[T, ID]) FindWhereExpr(query string, args ...interface{}) ([]T, error) {
	var models []T
	err := r.db.Where(query, args...).Find(&models).Error
	return models, err
}

// FindBy retrieves all records where the given field matches the value.
func (r *Repository[T, ID]) FindBy(field string, value any) ([]T, error) {
	var models []T
	err := r.db.Where(fmt.Sprintf("%s = ?", field), value).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find records: %w", err)
	}
	return models, nil
}

// FindFirstBy retrieves the first record where the field matches the value.
func (r *Repository[T, ID]) FindFirstBy(field string, value any) (T, error) {
	var model T
	err := r.db.Where(fmt.Sprintf("%s = ?", field), value).First(&model).Error
	if err != nil {
		return model, fmt.Errorf("failed to find first record: %w", err)
	}
	return model, nil
}

// FindById retrieves a record by its ID.
func (r *Repository[T, ID]) FindById(id ID) (T, error) {
	var model T
	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		return model, fmt.Errorf("failed to find record by ID: %w", err)
	}
	return model, nil
}

// FindByIds retrieves records by a list of IDs.
func (r *Repository[T, ID]) FindByIds(ids []ID) ([]T, error) {
	var models []T
	err := r.db.Where("id IN ?", ids).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find records by IDs: %w", err)
	}
	return models, nil
}

// Preload retrieves an entity and preloads the specified fields.
func (r *Repository[T, ID]) Preload(model *T, fields []string, key string, value any) error {
	db := r.db
	for _, field := range fields {
		db = db.Preload(field)
	}
	return db.Where(fmt.Sprintf("%s = ?", key), value).First(model).Error
}

// Exists checks if any record exists matching the query.
func (r *Repository[T, ID]) Exists(query string, args ...any) (bool, error) {
	var count int64
	err := r.db.Model(new(T)).Where(query, args...).Count(&count).Error
	return count > 0, err
}

// ExistByID checks if a record with the given ID exists.
func (r *Repository[T, ID]) ExistByID(id ID) (bool, error) {
	var count int64
	err := r.db.Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// SoftDelete marks an entity as deleted by setting a "deleted_at" timestamp.
func (r *Repository[T, ID]) SoftDelete(model T) error {
	return r.db.Model(model).Update("deleted_at", time.Now()).Error
}

func (r *Repository[T, ID]) SoftDeleteWithUpdate(model T, updates map[string]interface{}) error {
	// First, apply the additional updates like "updated_by"
	if len(updates) > 0 {
		if err := r.db.Model(model).Updates(updates).Error; err != nil {
			return err
		}
	}

	// Then set the deleted_at field to now
	return r.db.Model(model).Update("deleted_at", time.Now()).Error
}

// Delete deletes an entity from the database, with an option for hard delete.
func (r *Repository[T, ID]) Delete(model T, hard bool) error {
	var db = r.db
	if hard {
		db = db.Unscoped()
	}
	return db.Delete(model).Error
}

// DeleteById deletes a record by its ID, with an option for hard delete.
func (r *Repository[T, ID]) DeleteById(id ID, hard bool) error {
	var model T
	db := r.db
	if hard {
		db = db.Unscoped()
	}
	return db.Delete(&model, "id = ?", id).Error
}

// DeleteWhere deletes records matching the provided query.
func (r *Repository[T, ID]) DeleteWhere(query string, args ...any) error {
	return r.db.Where(query, args...).Delete(new(T)).Error
}

// DeleteAll deletes all records from the table.
func (r *Repository[T, ID]) DeleteAll(model T) error {
	return r.db.Where("1 = 1").Delete(&model).Error
}

// DeleteMany deletes multiple entities.
func (r *Repository[T, ID]) DeleteMany(models []T) error {
	return r.db.Delete(&models).Error
}

// DeleteManyByIds deletes records by a list of IDs.
func (r *Repository[T, ID]) DeleteManyByIds(model T, ids []ID) error {
	return r.db.Delete(&model, ids).Error
}

// Save saves the entity to the database.
func (r *Repository[T, ID]) Save(model T) (T, error) {
	err := r.db.Save(&model).Error
	if err != nil {
		return model, fmt.Errorf("failed to save entity: %w", err)
	}
	return model, nil
}

// SaveAll saves multiple records at once.
func (r *Repository[T, ID]) SaveAll(models []T) error {
	return r.db.Save(models).Error
}

// UpsertOnlyColumns updates the entity if it exists and creates it if it doesn't.
// conflictColumns specifies the columns to be used for conflict detection (i.e., the unique key or composite index).
// updateColumns specifies the columns to be updated in case of a conflict, ensuring only those columns are modified.
// The method performs an "upsert" operation, inserting a new record if no conflict is found, or updating the specified columns
// if a record with the same conflict columns (unique constraint or index) already exists.
func (r *Repository[T, ID]) UpsertOnlyColumns(model T, conflictColumns []string, updateColumns []string) (T, error) {
	// Dynamically construct conflict columns from the input
	var conflictClauseColumns []clause.Column
	for _, col := range conflictColumns {
		conflictClauseColumns = append(conflictClauseColumns, clause.Column{Name: col})
	}

	// Dynamically set the columns to be updated
	var updateClause = clause.AssignmentColumns(updateColumns)

	// Perform the upsert with specified conflict and update columns
	err := r.db.
		Clauses(clause.OnConflict{
			Columns:   conflictClauseColumns, // Specified conflict columns dynamically
			DoUpdates: updateClause,          // Only specified columns for update
		}).
		Create(model).Error

	if err != nil {
		return model, fmt.Errorf("failed to upsert entity with specified columns: %w", err)
	}
	return model, nil
}

// Updates performs a partial update on the given record using UpdateMap.
func (r *Repository[T, ID]) Updates(model T, m UpdateMap) error {

	if err := m.valid(); err != nil {
		return err
	}

	return r.db.Model(&model).Updates(m).Error
}

// clone creates a copy of the repository with the same DB connection.
func (r *Repository[T, ID]) clone() *Repository[T, ID] {
	return &Repository[T, ID]{db: r.db}
}
