package middleware

import (
	"github.com/gin-gonic/gin"
	"konsultn-api/internal/domain/team/service"
	"net/http"
)

func CanUpdateTeamMiddleware(teamService *service.TeamService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")

		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this team"})
			return
		}

		teamId := c.Param("id")

		allowed := teamService.CanUpdateOrDeleteTeam(teamId, userId.(string))

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this team"})
			return
		}
		c.Next()
	}
}
