package service

import (
	"fmt"
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/client"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/domain/team/repository"
	"konsultn-api/internal/shared/crud"
	"konsultn-api/internal/shared/helper"
)

type TeamService struct {
	teamRepo           *repository.TeamRepository
	teamMemberRepo     *repository.TeamMemberRepository
	teamInvitationRepo *repository.TeamInvitationRepository
	userClient         *client.UserClientImpl
	db                 *gorm.DB
}

func NewTeamService(db *gorm.DB) *TeamService {
	teamRepo := repository.NewTeamRepository(db)
	teamMemberRepo := repository.NewTeamMemberRepository(db)
	teamInvitationRepo := repository.NewTeamInvitationRepository(db)
	userRepo := crud.NewRepository[model.UserView](db) // You must import the user domain

	// Inject UserClientImpl with its dependencies
	userClient := &client.UserClientImpl{
		UserRepo: userRepo,
	}

	// Construct the TeamService
	return &TeamService{
		teamRepo:           teamRepo,
		teamMemberRepo:     teamMemberRepo,
		teamInvitationRepo: teamInvitationRepo,
		userClient:         userClient,
		db:                 db,
	}
}

func (s *TeamService) syncTeamMembers(teamId string, add []dto.AddMemberRequest, remove []string) error {
	for _, user := range add {
		// Skip if user doesn't exist in identity service
		if s.userClient.GetUserById(user.UserId).ID == "" {
			continue
		}

		member := model.TeamMember{
			TeamID: teamId,
			UserID: user.UserId,
			Role:   user.Role.String(),
		}

		// Try inserting, and if already exists, update only if the role has changed
		_, err := s.teamMemberRepo.UpsertOnlyColumns(&member, []string{"team_id", "user_id"}, []string{"role"})

		if err != nil {
			return fmt.Errorf("syncTeamMembers failed on user %s: %w", user.UserId, err)
		}
	}

	if len(remove) > 0 {
		if err := s.teamMemberRepo.DeleteWhere("team_id = ? AND user_id IN ?", teamId, remove); err != nil {
			return err
		}
	}

	return nil
}

func (s *TeamService) hydrateTeam(team *model.Team) error {
	// Preload team with members
	if err := s.teamRepo.Preload(team, []string{"Members"}, "id", team.ID); err != nil {
		return err
	}

	// Extract user IDs from team members
	memberIds := helper.Map(team.Members, func(m model.TeamMember) string {
		return m.UserID
	})

	// Fetch user data
	members := s.userClient.GetUsersByIds(memberIds)

	// Build a map of userID → UserView
	userMap := make(map[string]model.UserView, len(members))
	for _, u := range members {
		userMap[u.ID] = u
	}

	// Attach UserView to each member
	for i := range team.Members {
		if u, ok := userMap[team.Members[i].UserID]; ok {
			team.Members[i].User = u
		}
	}

	return nil
}
