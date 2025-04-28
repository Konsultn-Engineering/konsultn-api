package repository

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
)

type TeamMemberRepository[team model.TeamMember] struct {
	*crud.Repository[team]
}

func NewTeamMemberRepository(db *gorm.DB) *TeamMemberRepository[model.TeamMember] {
	return &TeamMemberRepository[model.TeamMember]{
		Repository: crud.NewRepository[model.TeamMember](db),
	}
}
