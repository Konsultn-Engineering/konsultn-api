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

	team := api.Group("/teams", middleware.AuthMiddleware())
	{
		team.GET("/:id", h.FindTeamById)
		team.POST("", h.CreateTeam)
		team.PATCH("/:id", middleware2.CanUpdateTeamMiddleware(teamService), h.UpdateTeamById)
		team.POST("/:id/invitations", middleware2.CanUpdateTeamMiddleware(teamService), h.InviteUsersToTeam)
	}

	teamInvitation := api.Group("/teams/invitations", middleware.AuthMiddleware())
	{
		teamInvitation.PATCH("/:invitationId/accept", h.AcceptInvitation)
		teamInvitation.PATCH("/:invitationId/reject", h.RejectInvitation)
	}
}
