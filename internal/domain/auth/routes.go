package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"konsultn-api/internal/domain/user"
)

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	auth := api.Group("/auth")
	repo := user.NewRepository(db)
	h := NewHandler(repo)
	{
		auth.POST("/register", h.CreateUser)
		auth.POST("/login", h.Login)
		//auth.POST("/refresh")
	}
}
