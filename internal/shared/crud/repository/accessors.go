package repository

import (
	"gorm.io/gorm"
)

func (r *BaseRepository[T, ID]) GetTableName() string {
	// Return cached table name if available
	if r.state.TableName != "" {
		return r.state.TableName
	}

	// Get the table name from GORM
	stmt := &gorm.Statement{DB: r.state.DB}
	err := stmt.Parse(new(T))

	if err != nil {
		return ""
	}

	// Cache the table name
	r.state.TableName = stmt.Schema.Table

	return r.state.TableName
}

func (r *BaseRepository[T, ID]) SetDB(db *gorm.DB) {
	r.state.DB = db
	r.db = db
}

func (r *BaseRepository[T, ID]) GetDB() *gorm.DB {
	return r.state.DB
}
