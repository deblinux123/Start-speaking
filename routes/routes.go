package routes

import (
	"start-speek/controllers"
	"start-speek/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(server *gin.Engine) {
	server.GET("/", home)
	server.POST("/register", controllers.Register)
	server.POST("/login", controllers.Login)

	authorized := server.Group("/api")

	authorized.Use(middleware.AuthMiddleWare())

	authorized.GET("/me", controllers.GetMe)
	authorized.PUT("/me", controllers.UpdateMe)
	authorized.PUT("/change-password", controllers.ChangePassword)
	authorized.POST("/logout", middleware.AuthMiddleWare(), controllers.LogOut)
	server.POST("/refresh", controllers.RefreshToken)
}
