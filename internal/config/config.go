package config

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	_ "konsultn-api/docs"
	"konsultn-api/internal/domain/auth"
	"konsultn-api/internal/domain/project"
	"konsultn-api/internal/domain/task"
	"konsultn-api/internal/domain/team"
	"konsultn-api/internal/domain/user"
)

func Setup(r *gin.Engine, db *gorm.DB) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := r.Group("/api")
	{
		user.RegisterRoutes(apiGroup, db)
		auth.RegisterRoutes(apiGroup, db)
		task.RegisterRoutes(apiGroup, db)
		project.RegisterRoutes(apiGroup, db)
		team.RegisterRoutes(apiGroup, db)
	}

}
