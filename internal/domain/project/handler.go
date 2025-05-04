package project

import (
	"github.com/gin-gonic/gin"
	"konsultn-api/internal/domain/project/dto"
	"konsultn-api/internal/domain/project/model"
	"konsultn-api/internal/domain/project/service"
	"konsultn-api/internal/shared/transport"
	"net/http"
)

type Handler struct {
	projectService *service.ProjectService
}

func NewHandler(projectService *service.ProjectService) *Handler {
	return &Handler{projectService: projectService}
}

// FindByID  godoc
// @Summary      Create a new project
// @Description  Creates a project with basic information
// @Tags         projects
// @Produce      json
// @Param        id        path      string  true   "Project ID"
// @Success      200 {object} model.Project
// @Router       /projects/{id} [get]
func (h *Handler) FindByID(ctx *gin.Context) {
	projectId := ctx.Param("id")
	var project model.Project

	h.projectService.Repo.Preload(&project, []string{"Tasks", "Tasks.Assignee"}, "id", projectId)

	ctx.JSON(http.StatusOK, project)
}

func (h *Handler) CreateProject(ctx *gin.Context) {
	var req model.Project
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, transport.ErrorResponse{Message: err.Error()})
		return
	}

	project, err := h.projectService.Repo.Save(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, project)
}

func (h *Handler) CreateProjectTask(ctx *gin.Context) {
	projectID := ctx.Param("id")

	var req dto.TaskDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, transport.ErrorResponse{Message: err.Error()})
		return
	}

	project, err := h.projectService.CreateProjectTask(projectID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, transport.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, project)
}
