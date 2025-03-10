package main

import (
	"gin-user-app/controllers"
	"gin-user-app/database"
	"gin-user-app/middleware"
	"log"

	_ "gin-user-app/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Your API
// @version 1.0
// @description Your API description
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /api

func main() {
	// Koneksi ke database
	database.ConnectDB()

	// Pastikan database berhasil terhubung
	if database.DB == nil {
		log.Fatal("❌ Database connection is nil. Exiting...")
	}

	// Setup router
	r := gin.Default()

	// Setup Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Group API Routes
	api := r.Group("/api")
	{
		// Routes untuk User (Panggil dari controllers)
		api.GET("/users", controllers.GetUsers)
		api.POST("/users", controllers.CreateUser)
		api.GET("/users/:id", controllers.GetUserByID)
		api.DELETE("/users/:id", controllers.DeleteUser)
		api.PUT("/users/:id", controllers.UpdateUser)

		// Routes untuk Auth
		api.POST("/auth/login", controllers.Login)
		api.GET("/auth/verify", middleware.AuthMiddleware(), controllers.VerifyUser)
		api.POST("/auth/logout", middleware.AuthMiddleware(), controllers.Logout)
	}

	// Jalankan server di port 8080
	r.Run(":8080")
}
