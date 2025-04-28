package task

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"konsultn-api/internal/middleware"
)

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	task := api.Group("/tasks", middleware.AuthMiddleware())
	repo := NewRepository(db)
	h := NewHandler(repo)
	{
		task.GET(":id", h.GetTaskById)
		task.POST("", h.CreateTask)
	}
}
