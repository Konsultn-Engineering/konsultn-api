package team

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"konsultn-api/internal/domain/team/service"
	"konsultn-api/internal/middleware"
)

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	teamService := service.NewTeamService(db)
	h := NewHandler(teamService)

	team := api.Group("/teams", middleware.AuthMiddleware())
	{
		team.POST("", h.CreateTeam)
	}
}
