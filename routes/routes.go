package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	api := server.Group("/api")
	api.GET("/games", getGames)
	api.POST("/games", createGame)

	api.POST("/publisher", createPublisher)

}
