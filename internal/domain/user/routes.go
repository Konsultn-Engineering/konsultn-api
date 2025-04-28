package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"konsultn-api/internal/middleware"
)

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	user := api.Group("/users", middleware.AuthMiddleware())
	repo := NewRepository(db)
	h := NewHandler(repo)
	{
		user.GET("", h.ListAllUsers)
		user.GET("/:id", h.GetUserById)
		user.POST("", h.CreateUser)
		user.DELETE("/:id", h.DeleteUser)
	}
}
