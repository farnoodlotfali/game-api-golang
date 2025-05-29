package main

import (
	"fmt"
	"net/http"

	"github.com/game-api/db"
	"github.com/game-api/objS3"
	"github.com/game-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Print("hi1")
	db.InitDB()
	objS3.InitS3()
	server := gin.Default()
	server.GET("/", func(context *gin.Context) { context.JSON(http.StatusOK, gin.H{"message": "Could not sag!"}) })

	routes.RegisterRoutes(server)

	server.Run(":8083")
}
