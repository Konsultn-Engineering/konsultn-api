package repository

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
)

type TeamQuery = crud.QueryBuilder[model.Team]
type TeamRepository struct {
	*crud.Repository[model.Team, string]
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{
		Repository: crud.NewRepository[model.Team, string](db),
	}
}

func (r TeamRepository) Members() TeamQuery {
	qb := r.Query()
	return qb.Join("team_members").On("id", "team_id")
}
