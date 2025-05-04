package repository

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/enum"
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
)

type TeamMemberRepository struct {
	*crud.Repository[model.TeamMember, string]
}

func NewTeamMemberRepository(db *gorm.DB) *TeamMemberRepository {
	return &TeamMemberRepository{
		Repository: crud.NewRepository[model.TeamMember, string](db),
	}
}

func (r *TeamMemberRepository) IsTeamAdmin(teamId string, userId string) bool {
	adminRoles := []string{enum.Owner.String(), enum.Admin.String()}

	records, err := r.FindWhereExpr(
		"team_id = ? AND user_id = ? AND role IN ?",
		teamId, userId, adminRoles,
	)

	if err != nil {
		// Optional: log or track this error
		return false
	}
	return len(records) > 0
}
