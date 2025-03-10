package routes

import (
	"gin-user-app/controllers"
	"gin-user-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/logout", controllers.Logout)
		auth.GET("/verify", middleware.AuthMiddleware(), controllers.VerifyUser)
	}
}
