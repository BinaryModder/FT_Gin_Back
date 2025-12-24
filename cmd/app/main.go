package main

import (
	"github.com/BinaryModder/FT_Gin_Back.git/internal/database"
	"github.com/BinaryModder/FT_Gin_Back.git/internal/handler"
	"github.com/BinaryModder/FT_Gin_Back.git/internal/repository"
	"github.com/BinaryModder/FT_Gin_Back.git/internal/service"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
)

func main() {
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.GET("", userHandler.GetAllUsers)
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	router.GET("/swagger/*any*", ginSwagger.WrapHandler(swaggerfiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
