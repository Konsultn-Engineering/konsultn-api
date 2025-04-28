package project

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"konsultn-api/internal/domain/project/service"
	"konsultn-api/internal/middleware"
)

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	projectService := service.NewProjectService(db)
	h := NewHandler(projectService)

	project := api.Group("/projects", middleware.AuthMiddleware())
	{
		project.GET("/:id", h.FindByID)
		project.POST("/", h.CreateProject)
		project.POST("/:id/tasks", h.CreateProjectTask)
	}
}
