package user

import (
	"gorm.io/gorm"
	"konsultn-api/internal/shared/crud"
)

type Repository[user User] struct {
	*crud.Repository[user, string]
}

func NewRepository(db *gorm.DB) *Repository[User] {
	return &Repository[User]{
		Repository: crud.NewRepository[User, string](db),
	}
}
