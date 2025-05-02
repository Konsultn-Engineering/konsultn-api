package repository

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
)

type TeamRepository struct {
	*crud.Repository[model.Team]
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{
		Repository: crud.NewRepository[model.Team](db),
	}
}
