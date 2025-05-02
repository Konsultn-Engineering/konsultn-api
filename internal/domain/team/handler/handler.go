package handler

import (
	service "konsultn-api/internal/domain/team/service"
)

type Handler struct {
	teamService *service.TeamService
}

func NewHandler(teamService *service.TeamService) *Handler {
	return &Handler{teamService: teamService}
}
