package task

import (
	"gorm.io/gorm"
	"konsultn-api/internal/shared/crud"
)

type Repository[task Task] struct {
	*crud.Repository[task]
}

func NewRepository(db *gorm.DB) *Repository[Task] {
	return &Repository[Task]{
		Repository: crud.NewRepository[Task](db),
	}
}
