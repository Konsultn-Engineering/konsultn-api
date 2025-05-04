package repository

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/project/model"
	"konsultn-api/internal/shared/crud"
)

type Repository[project model.Project] struct {
	*crud.Repository[project, string]
}

func NewRepository(db *gorm.DB) *Repository[model.Project] {
	return &Repository[model.Project]{
		Repository: crud.NewRepository[model.Project, string](db),
	}
}
