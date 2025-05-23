package types

import (
	"gorm.io/gorm"
)

// iCRUDRepository defines the core CRUD operations for a repository
type iCRUDRepository[T any, ID comparable] interface {
	SetDB(db *gorm.DB)
	GetDB() *gorm.DB
	GetTableName() string
	Select(fields []string) Repository[T, ID]
	Clone() Repository[T, ID]
	FindAll() ([]*T, error)
	Count() (int64, error)
	Query() QueryBuilder[T]
	FindWhere(filters map[string]interface{}) ([]*T, error)
	FindWhereExpr(query string, args ...interface{}) ([]*T, error)
	FindBy(field string, value any) ([]*T, error)
	FindFirstBy(field string, value any) (*T, error)
	FindById(id ID) (*T, error)
	FindByIds(ids []ID) ([]*T, error)
	Preload(model *T, fields []string, key string, value any) error
	Exists(query string, args ...any) (bool, error)
	ExistByID(id ID) (bool, error)
	Save(model *T) (*T, error)
	SaveAll(models []*T) error
	Updates(model *T, m UpdateMap) error
	UpsertOnlyColumns(model *T, conflictColumns []string, updateColumns []string) (*T, error)
	Delete(model *T, hard bool) error
	DeleteById(id ID, hard bool) error
	DeleteWhere(query string, args ...any) error
	DeleteAll() error
	DeleteMany(models []*T) error
	DeleteManyByIds(ids []ID) error
	SoftDelete(model *T) error
}

type Repository[T any, ID comparable] interface {
	iCRUDRepository[T, ID]
}
