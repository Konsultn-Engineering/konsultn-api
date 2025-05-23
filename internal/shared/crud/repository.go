package crud

import (
	"gorm.io/gorm"
	"konsultn-api/internal/shared/crud/repository"
)

// NewRepository creates a new BaseRepository instance with the provided database connection
// It returns a pointer to the newly created BaseRepository
func NewRepository[T any, ID comparable](db *gorm.DB) *Repository[T, ID] {
	return &Repository[T, ID]{
		BaseRepository: repository.NewBaseRepository[T, ID](db),
	}
}
