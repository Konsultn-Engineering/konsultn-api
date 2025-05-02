package handler

import (
	"github.com/gin-gonic/gin"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/shared/transport"
	"net/http"
)

func (h *Handler) CreateTeam(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var createTeamRequest dto.CreateTeamRequest
	if err := ctx.ShouldBindJSON(&createTeamRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	createdTeam, err := h.teamService.CreateTeam(userId.(string), createTeamRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToTeamDTO(&createdTeam))
}
func (h *Handler) FindTeamById(ctx *gin.Context) {
	teamId := ctx.Param("id")

	if teamId == "" {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: "invalid team id"})
		return
	}

	team, err := h.teamService.GetTeamById(teamId)

	if err != nil {
		ctx.JSON(http.StatusNotFound, transport.ErrorResponse{Message: "team not found"})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToTeamDTO(&team))
}
func (h *Handler) UpdateTeamById(ctx *gin.Context) {
	teamId := ctx.Param("id")
	userId, _ := ctx.Get("userId")

	var updateTeamRequest dto.UpdateTeamRequest
	if err := ctx.ShouldBindJSON(&updateTeamRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: "invalid update request"})
		return
	}

	updatedTeam, err := h.teamService.UpdateTeamById(userId.(string), teamId, updateTeamRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToTeamDTO(&updatedTeam))
}
