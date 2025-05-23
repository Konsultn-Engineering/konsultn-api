package handler

import (
	"github.com/gin-gonic/gin"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/shared/crud"
	"konsultn-api/internal/shared/transport"
	"net/http"
)

func (h *Handler) CreateTeam(ctx *gin.Context) {
	_, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var createTeamRequest dto.CreateTeamRequest
	if err := ctx.ShouldBindJSON(&createTeamRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	createdTeam, err := h.teamService.WithUser(ctx).CreateTeam(createTeamRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToTeamDTO(createdTeam))
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

	ctx.JSON(http.StatusOK, dto.ToTeamDTO(team))
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

	ctx.JSON(http.StatusOK, dto.ToTeamDTO(updatedTeam))
}

func (h *Handler) ListAllTeams(ctx *gin.Context) {
	var params crud.QueryParams

	// Bind basic query params
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	res := h.teamService.WithUser(ctx).Testing(params)

	ctx.JSON(http.StatusOK, res)

	///-----------------------------------------------
	//if filterMap, exists := ctx.Get("filterMap"); exists {
	//	params.Filter = filterMap.(map[string]string)
	//}
	//
	//teams, err := h.teamService.WithUser(ctx).GetAllUserTeams(params)
	//
	//if err != nil {
	//	ctx.JSON(http.StatusNotFound, "")
	//	return
	//}
	//
	//ctx.JSON(http.StatusOK, dto.ToTeamDTOPaginated(teams))
}
