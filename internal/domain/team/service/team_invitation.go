package service

import (
	"fmt"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/domain/team/enum"
	"konsultn-api/internal/domain/team/model"
	"time"
)

func (s *TeamService) InviteUsersToTeam(fromUserId string, teamId string, invitations []dto.AddMemberRequest) error {
	userIds := make([]string, 0, len(invitations))
	for _, inv := range invitations {
		userIds = append(userIds, inv.UserId)
	}

	validInvites, err := s.teamInvitationRepo.FindValidInvitations(teamId, userIds)
	if err != nil {
		return fmt.Errorf("fetching valid invitations: %w", err)
	}

	existing := make(map[string]*model.TeamInvitation, len(validInvites))
	for _, inv := range validInvites {
		existing[inv.ToUserID] = inv
	}

	for _, invitation := range invitations {
		if s.userClient.GetUserById(invitation.UserId).ID == "" {
			continue
		}

		if cur, ok := existing[invitation.UserId]; ok && cur.Role == invitation.Role.String() {
			continue
		}

		expiresAt := time.Now().Add(7 * 24 * time.Hour)
		newInvite := model.TeamInvitation{
			FromUserID: fromUserId,
			ToUserID:   invitation.UserId,
			TeamID:     teamId,
			Message:    invitation.Message,
			Status:     enum.Pending.String(),
			Role:       invitation.Role.String(),
			ExpiresAt:  &expiresAt,
		}

		if _, err := s.teamInvitationRepo.UpsertOnlyColumns(
			&newInvite,
			[]string{"team_id", "to_user_id"},
			[]string{"role", "status", "expires_at"},
		); err != nil {
			return fmt.Errorf("upserting invitation for user %s: %w", invitation.UserId, err)
		}

	}
	return nil
}

func (s *TeamService) UpdateTeamInvitation(invitationId string, action string, actingUserId string) error {
	// Fetch the invitation
	invitation, err := s.teamInvitationRepo.FindById(invitationId)
	if err != nil {
		return fmt.Errorf("error finding invitation: %w", err)
	}

	// Check if the invitation is for the acting user
	if invitation.ToUserID != actingUserId {
		return fmt.Errorf("unauthorized: you are not allowed to respond to this invitation")
	}

	// Check if the invitation is still valid
	if invitation.ExpiresAt != nil && invitation.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("invitation has expired")
	}

	// Only proceed if status is still pending
	if invitation.Status != enum.Pending.String() {
		return fmt.Errorf("invitation already responded to")
	}

	switch action {
	case "accept":
		invitation.Status = enum.Accepted.String()

		// Add the user as a team member
		_, err := s.teamMemberRepo.Save(&model.TeamMember{
			TeamID:   invitation.TeamID,
			UserID:   invitation.ToUserID,
			Role:     invitation.Role,
			JoinedAt: time.Now(),
		})
		if err != nil {
			return fmt.Errorf("error adding member to team: %w", err)
		}

	case "reject":
		invitation.Status = enum.Rejected.String()

	default:
		return fmt.Errorf("invalid action: must be 'accept' or 'reject'")
	}

	// Persist the updated invitation
	if _, err := s.teamInvitationRepo.Save(invitation); err != nil {
		return fmt.Errorf("failed to update invitation status: %w", err)
	}

	return nil
}
