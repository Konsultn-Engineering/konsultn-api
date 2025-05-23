package builder

import (
	"gorm.io/gorm"
	"konsultn-api/internal/shared/crud/builder/utils"
	"konsultn-api/internal/shared/crud/types"
)

// QueryBuilder main struct
type QueryBuilder[T any] struct {
	*gorm.DB
	model           *T
	baseTable       string
	lastJoinedTable string
	knownAliases    map[string]bool
	orderClause     string
	joins           []JoinClause
	page            int
	limit           int
}

// NewQueryBuilder creates a new query builder for the given model
func NewQueryBuilder[T any](db *gorm.DB) types.QueryBuilder[T] {
	var model T

	return &QueryBuilder[T]{
		DB:           db,
		model:        &model,
		baseTable:    getTableName[T](db),
		knownAliases: make(map[string]bool),
		page:         1,
		limit:        10,
	}
}

func getTableName[T any](db *gorm.DB) string {
	var model T
	stmt := &gorm.Statement{DB: db}
	_ = stmt.Parse(&model)
	return stmt.Schema.Table
}

// G returns the underlying gorm.DB with or without scope
func (qb *QueryBuilder[T]) G(scoped bool) *gorm.DB {
	if !scoped {
		return qb.build().Unscoped()
	}
	return qb.build()
}

// Unscoped returns a new query builder without model's default scopes
func (qb *QueryBuilder[T]) Unscoped() types.QueryBuilder[T] {
	qb.DB = qb.DB.Unscoped()
	return qb
}

// build assembles the query with all conditions and joins
func (qb *QueryBuilder[T]) build() *gorm.DB {
	db := qb.DB.Session(&gorm.Session{NewDB: true})
	db = qb.buildJoins()
	return db
}

// Raw creates a safe raw SQL fragment
// For fixed SQL expressions (no parameters):
//
//	qb.Raw("NOW()")
//
// For parameterized SQL:
//
//	qb.Raw("EXTRACT(YEAR FROM ?)", someDate)
func (qb *QueryBuilder[T]) Raw(sql string, args ...interface{}) types.RawValue {
	return utils.SafeSQL(sql, args...)
}

// ToRawSQL converts the query to its SQL representation
func (qb *QueryBuilder[T]) ToRawSQL() string {
	var models T
	sql := qb.build().ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(&models)
	})
	return sql
}
