package crud

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// Repository is a generic repository that provides CRUD methods for any type.
type Repository[T any] struct {
	db           *gorm.DB
	selectFields []string
}

// NewRepository creates and returns a new instance of Repository[T].
func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

// Select sets the fields to be selected in the query.
func (r *Repository[T]) Select(fields []string) *Repository[T] {
	newRepo := r.clone()
	if len(fields) > 0 {
		newRepo.db = newRepo.db.Select(fields)
	}
	return newRepo
}

// FindAll retrieves all records for the given model type.
func (r *Repository[T]) FindAll() ([]T, error) {
	var models []T
	err := r.db.Find(&models).Error
	return models, err
}

// FindWhere retrieves records from the database where all provided field-value pairs match exactly.
// Useful for filtering by multiple columns with AND conditions.
// Example: FindWhere(map[string]interface{}{"team_id": "abc", "role": "owner"})
func (r *Repository[T]) FindWhere(filters map[string]interface{}) ([]T, error) {
	var models []T
	err := r.db.Where(filters).Find(&models).Error
	return models, err
}

// FindWhereExpr retrieves records using a custom WHERE clause and arguments.
// Allows more complex queries (e.g., OR conditions, IN clauses, LIKE, etc.).
// Example: FindWhereExpr("status = ? OR priority = ?", "active", "high")
func (r *Repository[T]) FindWhereExpr(query string, args ...interface{}) ([]T, error) {
	var models []T
	err := r.db.Where(query, args...).Find(&models).Error
	return models, err
}

// FindBy retrieves all records where the given field matches the value.
func (r *Repository[T]) FindBy(field string, value any) ([]T, error) {
	var models []T
	err := r.db.Where(fmt.Sprintf("%s = ?", field), value).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find records: %w", err)
	}
	return models, nil
}

// FindFirstBy retrieves the first record where the field matches the value.
func (r *Repository[T]) FindFirstBy(field string, value any) (*T, error) {
	var model T
	err := r.db.Where(fmt.Sprintf("%s = ?", field), value).First(&model).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find first record: %w", err)
	}
	return &model, nil
}

// FindById retrieves a record by its ID.
func (r *Repository[T]) FindById(id string) (*T, error) {
	var model T
	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find record by ID: %w", err)
	}
	return &model, nil
}

// FindByIds retrieves records by a list of IDs.
func (r *Repository[T]) FindByIds(ids []string) ([]T, error) {
	var models []T
	err := r.db.Where("id IN ?", ids).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find records by IDs: %w", err)
	}
	return models, nil
}

// Preload retrieves an entity and preloads the specified fields.
func (r *Repository[T]) Preload(entity *T, fields []string, key string, value any) error {
	db := r.db
	for _, field := range fields {
		db = db.Preload(field)
	}
	return db.Where(fmt.Sprintf("%s = ?", key), value).First(entity).Error
}

// Save saves the entity to the database (either creates or updates it).
func (r *Repository[T]) Save(entity *T) (*T, error) {
	err := r.db.Save(&entity).Error
	if err != nil {
		return nil, fmt.Errorf("failed to save entity: %w", err)
	}
	return entity, nil
}

// UpsertOnlyColumns updates the entity if it exists and creates it if it doesn't.
// conflictColumns specifies the columns to be used for conflict detection (i.e., the unique key or composite index).
// updateColumns specifies the columns to be updated in case of a conflict, ensuring only those columns are modified.
// The method performs an "upsert" operation, inserting a new record if no conflict is found, or updating the specified columns
// if a record with the same conflict columns (unique constraint or index) already exists.
func (r *Repository[T]) UpsertOnlyColumns(entity *T, conflictColumns []string, updateColumns []string) (*T, error) {
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
		Create(entity).Error

	if err != nil {
		return nil, fmt.Errorf("failed to upsert entity with specified columns: %w", err)
	}
	return entity, nil
}

func (r *Repository[T]) Exists(query string, args ...any) (bool, error) {
	var count int64
	err := r.db.Model(new(T)).Where(query, args...).Count(&count).Error
	return count > 0, err
}

// SoftDelete marks an entity as deleted by setting a "deleted_at" timestamp.
func (r *Repository[T]) SoftDelete(entity *T) error {
	return r.db.Model(entity).Update("deleted_at", time.Now()).Error
}

// Delete deletes an entity from the database, with an option for hard delete.
func (r *Repository[T]) Delete(entity *T, hard bool) error {
	var db = r.db
	if hard {
		db = db.Unscoped()
	}
	return db.Delete(entity).Error
}

// DeleteById deletes a record by its ID, with an option for hard delete.
func (r *Repository[T]) DeleteById(id string, hard bool) error {
	var model T
	db := r.db
	if hard {
		db = db.Unscoped()
	}
	return db.Delete(&model, "id = ?", id).Error
}

func (r *Repository[T]) DeleteWhere(query string, args ...any) error {
	return r.db.Where(query, args...).Delete(new(T)).Error
}

// clone creates a copy of the repository with the same DB connection.
func (r *Repository[T]) clone() *Repository[T] {
	return &Repository[T]{db: r.db}
}
