package repository

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
)

type TeamRepository[team model.Team] struct {
	*crud.Repository[team]
}

func NewTeamRepository(db *gorm.DB) *TeamRepository[model.Team] {
	return &TeamRepository[model.Team]{
		Repository: crud.NewRepository[model.Team](db),
	}
}
