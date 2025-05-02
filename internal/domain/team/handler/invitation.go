package handler

import (
	"github.com/gin-gonic/gin"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/shared/transport"
	"net/http"
)

func (h *Handler) InviteUsersToTeam(ctx *gin.Context) {
	teamId := ctx.Param("id")
	fromUserId, _ := ctx.Get("userId")

	var addMemberRequest []dto.AddMemberRequest

	if err := ctx.ShouldBindJSON(&addMemberRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	err := h.teamService.InviteUsersToTeam(fromUserId.(string), teamId, addMemberRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "success"})
}
func (h *Handler) AcceptInvitation(ctx *gin.Context) {
	invitationId := ctx.Param("invitationId")
	actingUserId, _ := ctx.Get("userId")

	if err := h.teamService.UpdateTeamInvitation(invitationId, "accept", actingUserId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "invitation accepted successfully"})
}

func (h *Handler) RejectInvitation(ctx *gin.Context) {
	invitationId := ctx.Param("invitationId")
	actingUserId, _ := ctx.Get("userId")

	if err := h.teamService.UpdateTeamInvitation(invitationId, "reject", actingUserId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "invitation rejected successfully"})
}
