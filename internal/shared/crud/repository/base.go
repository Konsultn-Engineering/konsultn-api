package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"konsultn-api/internal/shared/crud/builder"
	"konsultn-api/internal/shared/crud/types"
	"time"
)

type BaseRepository[T any, ID comparable] struct {
	state *types.RepositoryState[T, ID]
	db    *gorm.DB
}

// NewBaseRepository creates a new BaseRepository instance with the provided database connection
// It returns a pointer to the newly created BaseRepository
func NewBaseRepository[T any, ID comparable](db *gorm.DB) *BaseRepository[T, ID] {
	state := types.NewRepositoryState[T, ID](db)
	return &BaseRepository[T, ID]{
		state: state,
		db:    state.DB,
	}
}

// Verify implementation at compile time
var _ types.Repository[any, string] = (*BaseRepository[any, string])(nil)

// Clone creates a copy of the current BaseRepository instance
// It returns a new BaseRepository instance with the same database connection
func (r *BaseRepository[T, ID]) Clone() types.Repository[T, ID] {
	newState := r.state.Clone()
	return &BaseRepository[T, ID]{
		state: newState,
		db:    newState.DB,
	}
}

func (r *BaseRepository[T, ID]) WithContext(ctx context.Context) types.Repository[T, ID] {
	repo := r.Clone()
	repo.SetDB(repo.GetDB().WithContext(ctx))
	return repo
}

// Select creates a new BaseRepository instance with specified fields to be selected in queries
// It returns a new BaseRepository instance with the modified database query
func (r *BaseRepository[T, ID]) Select(fields []string) types.Repository[T, ID] {
	newRepo := r.Clone()
	if len(fields) > 0 {
		newRepo.SetDB(newRepo.GetDB().Select(fields))
	}
	return newRepo
}

// FindAll retrieves all records from the table
// It returns a slice of model pointers and any error encountered
func (r *BaseRepository[T, ID]) FindAll() ([]*T, error) {
	var models []*T
	err := r.db.Find(&models).Error
	return models, err
}

// Count returns the total number of records in the table
// It returns the count as int64 and any error encountered
func (r *BaseRepository[T, ID]) Count() (int64, error) {
	var count int64
	err := r.db.Model(new(T)).Count(&count).Error
	return count, err
}

// Query executes an advanced query with support for joins, where clauses, filtering, pagination, and sorting
func (r *BaseRepository[T, ID]) Query() types.QueryBuilder[T] {
	dbClone := r.db.Session(&gorm.Session{}) // clone without inherited settings
	// optionally, set model to T to scope queries
	dbClone = dbClone.Model(new(T))

	return builder.NewQueryBuilder[T](dbClone)
}

// FindWhere retrieves records matching the provided filter conditions
// It returns a slice of model pointers and any error encountered
func (r *BaseRepository[T, ID]) FindWhere(filters map[string]interface{}) ([]*T, error) {
	var models []*T
	err := r.db.Where(filters).Find(&models).Error
	return models, err
}

// FindWhereExpr retrieves records matching the provided query expression and arguments
// It returns a slice of model pointers and any error encountered
func (r *BaseRepository[T, ID]) FindWhereExpr(query string, args ...interface{}) ([]*T, error) {
	var models []*T
	err := r.db.Where(query, args...).Find(&models).Error
	return models, err
}

// FindBy retrieves records where the specified field matches the provided value
// It returns a slice of model pointers and any error encountered
func (r *BaseRepository[T, ID]) FindBy(field string, value any) ([]*T, error) {
	var models []*T
	err := r.db.Where(fmt.Sprintf("%s = ?", field), value).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find records: %w", err)
	}
	return models, err
}

// FindFirstBy retrieves the first record where the specified field matches the provided value
// It returns a single model pointer and any error encountered
func (r *BaseRepository[T, ID]) FindFirstBy(field string, value any) (*T, error) {
	var model T
	err := r.db.Where(fmt.Sprintf("%s = ?", field), value).First(&model).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find first record: %w", err)
	}
	return &model, nil
}

// FindById retrieves a record by its ID
// It returns a single model pointer and any error encountered
func (r *BaseRepository[T, ID]) FindById(id ID) (*T, error) {
	var model T
	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find record by ID: %w", err)
	}
	return &model, nil
}

// FindByIds retrieves multiple records by their IDs
// It returns a slice of model pointers and any error encountered
func (r *BaseRepository[T, ID]) FindByIds(ids []ID) ([]*T, error) {
	var models []*T
	err := r.db.Where("id IN ?", ids).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find records by IDs: %w", err)
	}
	return models, err
}

// Preload loads the specified relationships for a model with a condition
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) Preload(model *T, fields []string, key string, value any) error {
	db := r.db
	for _, field := range fields {
		db = db.Preload(field)
	}
	return db.Where(fmt.Sprintf("%s = ?", key), value).First(model).Error
}

// Exists checks if any record exists matching the provided query and arguments
// It returns a boolean indicating existence and any error encountered
func (r *BaseRepository[T, ID]) Exists(query string, args ...any) (bool, error) {
	var count int64
	err := r.db.Model(new(T)).Where(query, args...).Count(&count).Error
	return count > 0, err
}

// ExistByID checks if a record exists with the specified ID
// It returns a boolean indicating existence and any error encountered
func (r *BaseRepository[T, ID]) ExistByID(id ID) (bool, error) {
	var count int64
	err := r.db.Model(new(T)).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Save creates or updates a model in the database
// It returns the saved model and any error encountered
func (r *BaseRepository[T, ID]) Save(model *T) (*T, error) {
	err := r.db.Save(model).Error
	if err != nil {
		return nil, fmt.Errorf("failed to save entity: %w", err)
	}
	return model, nil
}

// SaveAll saves multiple models in a single operation
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) SaveAll(models []*T) error {
	if len(models) == 0 {
		return nil
	}
	return r.db.Save(&models).Error
}

// Updates specific fields of a model using the provided update map
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) Updates(model *T, m types.UpdateMap) error {
	if err := m.Valid(); err != nil {
		return err
	}
	return r.db.Model(model).Updates(m).Error
}

// UpsertOnlyColumns performs an upsert operation with specified conflict and update columns
// It returns the upserted model and any error encountered
func (r *BaseRepository[T, ID]) UpsertOnlyColumns(model *T, conflictColumns []string, updateColumns []string) (*T, error) {
	var conflictClauseColumns []clause.Column
	for _, col := range conflictColumns {
		conflictClauseColumns = append(conflictClauseColumns, clause.Column{Name: col})
	}

	updateClause := clause.AssignmentColumns(updateColumns)

	err := r.db.Clauses(clause.OnConflict{
		Columns:   conflictClauseColumns,
		DoUpdates: updateClause,
	}).Create(model).Error

	if err != nil {
		return nil, fmt.Errorf("failed to upsert entity with specified columns: %w", err)
	}
	return model, nil
}

// Delete performs either a soft or hard delete on the model based on the hard parameter
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) Delete(model *T, hard bool) error {
	db := r.db
	if hard {
		db = db.Unscoped()
	}
	return db.Delete(model).Error
}

// DeleteById deletes a record by its ID, either soft or hard delete based on the hard parameter
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) DeleteById(id ID, hard bool) error {
	db := r.db
	if hard {
		db = db.Unscoped()
	}
	return db.Delete(new(T), "id = ?", id).Error
}

// DeleteWhere deletes records matching the provided query and arguments
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) DeleteWhere(query string, args ...any) error {
	return r.db.Where(query, args...).Delete(new(T)).Error
}

// DeleteAll deletes all records from the table
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(new(T)).Error
}

// DeleteMany deletes multiple models in a single operation
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) DeleteMany(models []*T) error {
	if len(models) == 0 {
		return nil
	}
	return r.db.Delete(&models).Error
}

// DeleteManyByIds deletes multiple records by their IDs
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) DeleteManyByIds(ids []ID) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Delete(new(T), ids).Error
}

// SoftDelete performs a soft delete on the model by setting its deleted_at timestamp
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) SoftDelete(model *T) error {
	return r.db.Model(model).Update("deleted_at", time.Now()).Error
}

// SoftDeleteWithUpdate performs a soft delete with additional updates to the model
// It returns any error encountered during the operation
func (r *BaseRepository[T, ID]) SoftDeleteWithUpdate(model *T, updates map[string]interface{}) error {
	if len(updates) > 0 {
		if err := r.db.Model(model).Updates(updates).Error; err != nil {
			return err
		}
	}
	return r.db.Model(model).Update("deleted_at", time.Now()).Error
}
