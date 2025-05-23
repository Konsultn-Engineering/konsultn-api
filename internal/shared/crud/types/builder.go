package types

import (
	"gorm.io/gorm"
	"konsultn-api/internal/shared/crud/pagination"
)

type RawValue struct {
	Value string
	Args  []interface{}
}

type JoinBuilder interface {
	Raw(value interface{}) RawValue
	RawSQL(sql string, values ...interface{}) RawValue
	On(left, op, right interface{})
	And(left, op, right interface{})
	Or(left, op, right interface{})
	String() string
	GetParams() []interface{}
}

type QueryBuilder[T any] interface {
	// G Core/Base Methods
	G(scoped bool) *gorm.DB
	Unscoped() QueryBuilder[T]
	ToRawSQL() string

	// Select Methods
	Select(fields ...interface{}) QueryBuilder[T]

	// Where Conditions
	Where(field string, value interface{}) QueryBuilder[T]
	WhereNot(field string, value interface{}) QueryBuilder[T]
	WhereLT(field string, value interface{}) QueryBuilder[T]
	WhereLTE(field string, value interface{}) QueryBuilder[T]
	WhereGT(field string, value interface{}) QueryBuilder[T]
	WhereGTE(field string, value interface{}) QueryBuilder[T]
	WhereIN(field string, values interface{}) QueryBuilder[T]
	WhereNotIN(field string, values interface{}) QueryBuilder[T]
	WhereBetween(field string, min, max interface{}) QueryBuilder[T]
	WhereRaw(sql string, args ...interface{}) QueryBuilder[T]
	WhereGroup(callback func(QueryBuilder[T])) QueryBuilder[T]

	// OrWhere Conditions
	OrWhere(field string, value interface{}) QueryBuilder[T]
	OrWhereNot(field string, value interface{}) QueryBuilder[T]
	OrWhereGTE(field string, value interface{}) QueryBuilder[T]
	OrWhereLTE(field string, value interface{}) QueryBuilder[T]
	OrWhereIN(field string, value interface{}) QueryBuilder[T]
	OrWhereNotIN(field string, value interface{}) QueryBuilder[T]
	OrWhereBetween(field string, min, max interface{}) QueryBuilder[T]
	OrWhereRaw(sql string, args ...interface{}) QueryBuilder[T]
	OrWhereGroup(callback func(QueryBuilder[T])) QueryBuilder[T]

	// Join Operations
	Join(table string, opts ...string) QueryBuilder[T]
	LeftJoin(table string, opts ...string) QueryBuilder[T]
	RightJoin(table string, opts ...string) QueryBuilder[T]
	CrossJoin(table string, opts ...string) QueryBuilder[T]
	On(left, right string) QueryBuilder[T]
	OnGroup(builder func(joinBuilder JoinBuilder)) QueryBuilder[T]

	//RawSelect +
	RawSelect(expression string, alias string, args ...interface{}) QueryBuilder[T]
	Now() RawValue
	Cast(value interface{}, dataType string) RawValue
	Coalesce(fields []string, defaultValue ...interface{}) RawValue

	// GroupBy and Having Clauses
	GroupBy(fields ...string) QueryBuilder[T]
	Having(condition string, args ...interface{}) QueryBuilder[T]
	HavingEQ(field string, value interface{}) QueryBuilder[T]
	HavingNEQ(field string, value interface{}) QueryBuilder[T]
	HavingGT(field string, value interface{}) QueryBuilder[T]
	HavingGTE(field string, value interface{}) QueryBuilder[T]
	HavingLT(field string, value interface{}) QueryBuilder[T]
	HavingLTE(field string, value interface{}) QueryBuilder[T]
	HavingIN(field string, values interface{}) QueryBuilder[T]
	HavingBetween(field string, min, max interface{}) QueryBuilder[T]
	OrHaving(condition string, args ...interface{}) QueryBuilder[T]
	HavingGroup(callback func(QueryBuilder[T])) QueryBuilder[T]
	OrHavingGroup(callback func(QueryBuilder[T])) QueryBuilder[T]

	// WithPageParams And Paginate Pagination
	WithPageParams(params pagination.QueryParams) QueryBuilder[T]
	Paginate() (*pagination.PaginatedResult[T], error)
	PaginateMap() (pagination.PaginatedResult[map[string]interface{}], error)

	//First etc. Fetch Retrieval Methods
	First() (*T, error)
	FirstAsMap() (map[string]interface{}, error)
	All() ([]T, error)
	AllAsMaps() ([]map[string]interface{}, error)
	Into(dest interface{}) error

	Count() (int64, error)
	Exists() (bool, error)
}
