package team

import (
	"errors"
	"github.com/gin-gonic/gin"
	"konsultn-api/internal/domain/team/dto"
	service "konsultn-api/internal/domain/team/service"
	"konsultn-api/internal/shared/transport"
	"net/http"
)

type Handler struct {
	teamService *service.TeamService
}

func NewHandler(teamService *service.TeamService) *Handler {
	return &Handler{teamService: teamService}
}

func GetUserIDFromContext(c *gin.Context) (string, error) {
	userIDValue, exists := c.Get("userId")
	if !exists {
		return "", errors.New("userId not found in context")
	}

	userID, ok := userIDValue.(string)
	if !ok {
		return "", errors.New("userId in context is not a string")
	}

	return userID, nil
}

func (h *Handler) CreateTeam(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var createTeamRequest dto.CreateTeamRequest
	if err := ctx.ShouldBindJSON(&createTeamRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: "here"})
		return
	}

	createdTeam, err := h.teamService.CreateTeam(userId.(string), createTeamRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdTeam)
}
