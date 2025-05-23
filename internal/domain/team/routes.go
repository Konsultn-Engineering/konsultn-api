package team

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/handler"
	middleware2 "konsultn-api/internal/domain/team/middleware"
	"konsultn-api/internal/domain/team/service"
	"konsultn-api/internal/middleware"
)

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	teamService := service.NewTeamService(db)
	h := handler.NewHandler(teamService)
	canUpdateTeamMiddleware := middleware2.CanUpdateTeam(teamService)

	teams := api.Group("/teams", middleware.AuthMiddleware())
	{
		// Basic team operations
		teams.POST("", h.CreateTeam)      // Create a team
		teams.GET("/:id", h.FindTeamById) // Get team by ID
		teams.GET("", middleware.FilterMapMiddleware(), h.ListAllTeams)
		teams.PATCH("/:id", canUpdateTeamMiddleware, h.UpdateTeamById) // Update team metadata

		// Member management under a team
		members := teams.Group("/:id/members")
		{
			members.DELETE("/:memberId", canUpdateTeamMiddleware, h.RemoveTeamMember)    // Remove a member
			members.PATCH("/:memberId", canUpdateTeamMiddleware, h.UpdateTeamMemberById) // (Optional) Update member role, etc.
		}

		// Invitations under a team
		teamInvitations := teams.Group("/:id/invitations", canUpdateTeamMiddleware)
		{
			teamInvitations.POST("", h.InviteUsersToTeam) // Invite users to team
		}

		// Global invitation actions (accept/reject)
		invitations := teams.Group("/invitations")
		{
			invitations.PATCH("/:invitationId/accept", h.AcceptInvitation)  // Accept invite
			invitations.DELETE("/:invitationId/reject", h.RejectInvitation) // Reject invite
		}
	}
}
