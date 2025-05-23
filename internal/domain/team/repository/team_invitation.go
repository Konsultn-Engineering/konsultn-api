package repository

import (
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/shared/crud"
	"time"
)

type TeamInvitationRepository struct {
	*crud.Repository[model.TeamInvitation, string]
}

func NewTeamInvitationRepository(db *gorm.DB) *TeamInvitationRepository {
	return &TeamInvitationRepository{
		Repository: crud.NewRepository[model.TeamInvitation, string](db),
	}
}

func (r *TeamInvitationRepository) FindValidInvitations(teamID string, userIds []string) ([]*model.TeamInvitation, error) {
	var invitations []*model.TeamInvitation

	invitations, err := r.FindWhereExpr("team_id = ? AND to_user_id IN ? AND expires_at > ?", teamID, userIds, time.Now())

	if err != nil {
		return nil, err
	}
	return invitations, nil
}
