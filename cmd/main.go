package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"konsultn-api/internal/config"
	"konsultn-api/internal/db"
	"konsultn-api/internal/domain/project/model"
	"konsultn-api/internal/domain/task"
	model2 "konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/domain/user"
	"konsultn-api/pkg/firebase"
	"log"
	"os"
	"time"
)

func EnableGormSQLLogging(db *gorm.DB) {
	db.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, Konsultn! Coming soon!"})
	})

	firebaseErr := firebase.InitFirebase()
	if firebaseErr != nil {
		print(firebaseErr.Error())
		return
	}

	connection, _ := db.InitDB()
	//connection = connection.Debug()
	errr := connection.AutoMigrate(&user.User{}, &model.Project{}, &task.Task{}, model2.Team{}, model2.TeamMember{}, model2.TeamInvitation{})
	if errr != nil {
		return
	}

	config.Setup(r, connection)

	/*
	* Handle port and host from environment
	 */

	if gin.Mode() == gin.DebugMode {
		EnableGormSQLLogging(connection)
	}

	host := os.Getenv("IP") // Use "IP" for the host
	if host == "" {
		host = "::" // Default to all types if not set
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8100" // fallback
	}
	err := r.Run("[" + host + "]" + ":" + port)

	if err != nil {
		return
	} // default is :8080
}
