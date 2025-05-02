package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"konsultn-api/internal/config"
	"konsultn-api/internal/db"
	"konsultn-api/internal/domain/project/model"
	"konsultn-api/internal/domain/task"
	model2 "konsultn-api/internal/domain/team/model"
	"konsultn-api/internal/domain/user"
	"konsultn-api/pkg/firebase"
	"os"
)

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

	connection.Create(&user.User{Email: "khalid.elokiely@gmail.com"})
	connection.Find(&user.User{})

	var users []user.User
	result := connection.Find(&users)

	fmt.Println("Result error:", result.Error) // Capture any error
	fmt.Printf("Users: %+v\n", users)          // Log the results

	/*
	* Handle port and host from environment
	 */

	host := os.Getenv("IP") // Use "IP" for the host
	if host == "" {
		host = "::" // Default to all interfaces if not set
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
