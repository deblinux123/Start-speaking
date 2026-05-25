package main

import (
	"start-speek/db"
	"start-speek/models"
	"start-speek/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	db.InitRedis()

	db.DB.AutoMigrate(&models.User{})

	server := gin.Default()

	routes.RegisterRouter(server)
	server.Run(":8080")
}
