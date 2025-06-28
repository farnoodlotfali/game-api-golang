package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/game-api/db"
	"github.com/game-api/objS3"
	"github.com/game-api/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting server...")

	// Initialize database and S3
	db.InitDB()
	objS3.InitS3()

	// Create a Gin router with default middleware (logger and recovery)
	server := gin.Default()

	// CORS configuration
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Change "*" to your frontend domain in production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Root endpoint
	server.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Server is running!"})
	})

	// Register application routes
	routes.RegisterRoutes(server)

	// Start server on port 8083
	if err := server.Run(":8083"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
