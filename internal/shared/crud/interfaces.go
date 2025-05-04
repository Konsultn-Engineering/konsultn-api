package crud

type iCRUDRepository[T any, ID idType] interface {
	Select(fields []string) *Repository[T, ID]
	Count() (int64, error)
	FindAll() ([]T, error)
	FindWhere(filters map[string]interface{}) ([]T, error)
	FindWhereExpr(query string, args ...interface{}) ([]T, error)
	FindBy(field string, value any) ([]T, error)
	FindFirstBy(field string, value any) (*T, error)
	FindById(id ID) (*T, error)
	FindByIds(ids []ID) ([]T, error)
	Preload(entity *T, fields []string, key string, value any) error
	UpsertOnlyColumns(entity *T, conflictColumns []string, updateColumns []string) (*T, error)
	Exists(query string, args ...any) (bool, error)
	ExistByID(id ID) (bool, error)
	SoftDelete(entity *T) error
	Delete(entity *T, hard bool) error
	DeleteById(id ID, hard bool) error
	DeleteWhere(query string, args ...any) error
	DeleteAll(model T) error
	DeleteMany(models []T) error
	DeleteManyByIds(model T, ids []ID) error
	Save(entity *T) (*T, error)
	SaveAll(ts []T) error
	Updates(t T, m UpdateMap) error
}
