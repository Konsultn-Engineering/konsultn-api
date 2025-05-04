package handler

import (
	"github.com/gin-gonic/gin"
	"konsultn-api/internal/domain/team/dto"
	"konsultn-api/internal/shared/transport"
	"net/http"
)

func (h *Handler) UpdateTeamMemberById(ctx *gin.Context) {
	teamId := ctx.Param("id")
	memberId := ctx.Param("memberId")

	var updateMemberRequest dto.UpdateMemberRequest

	if err := ctx.ShouldBindJSON(&updateMemberRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	err := h.teamService.WithUser(ctx).UpdateTeamMember(teamId, memberId, updateMemberRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "update successful"})
}

func (h *Handler) UpdateTeamMembers(ctx *gin.Context) {

}

func (h *Handler) RemoveTeamMember(ctx *gin.Context) {
	teamId := ctx.Param("id")
	memberId := ctx.Param("memberId")

	err := h.teamService.WithUser(ctx).RemoveTeamMember(teamId, memberId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "member removed"})
}
