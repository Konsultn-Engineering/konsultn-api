package types

import (
	"gorm.io/gorm"
)

// RepositoryState encapsulates the common state needed by repositories
type RepositoryState[T any, ID comparable] struct {
	DB           *gorm.DB
	SelectFields []string
	TableName    string
}

// NewRepositoryState creates a new state instance
func NewRepositoryState[T any, ID comparable](db *gorm.DB) *RepositoryState[T, ID] {
	return &RepositoryState[T, ID]{
		DB: db,
	}
}

// Clone creates a copy of the state
func (s *RepositoryState[T, ID]) Clone() *RepositoryState[T, ID] {
	return &RepositoryState[T, ID]{
		DB:           s.DB,
		SelectFields: append([]string{}, s.SelectFields...),
		TableName:    s.TableName,
	}
}
