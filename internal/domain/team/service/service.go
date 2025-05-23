package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/client"
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
	actingUserId       string
}

func (s *TeamService) WithUser(ctx *gin.Context) *TeamService {
	// Return a shallow copy with user context
	userId, _ := ctx.Get("userId")
	newTeamService := *s // creates a shallow copy of the struct
	newTeamService.actingUserId = userId.(string)
	return &newTeamService
}

func NewTeamService(db *gorm.DB) *TeamService {
	teamRepo := repository.NewTeamRepository(db)
	teamMemberRepo := repository.NewTeamMemberRepository(db)
	teamInvitationRepo := repository.NewTeamInvitationRepository(db)
	userRepo := crud.NewRepository[model.UserView, string](db)

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

func (s *TeamService) hydrateOwner(team *model.Team) error {
	owner := s.userClient.GetUserById(team.OwnerID)
	team.Owner = &owner
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
	userMap := make(map[string]*model.UserView, len(members))
	for _, u := range members {
		userMap[u.ID] = u
	}

	// Attach UserView to each member
	for i := range team.Members {
		if u, ok := userMap[team.Members[i].UserID]; ok {
			team.Members[i].User = *u
		}
	}

	return nil
}
