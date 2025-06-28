package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	api := server.Group("/api")

	// games
	api.GET("/games", getGames)
	api.POST("/games", createGame)
	api.PUT("/games/:id", updateGame)
	api.GET("/games/:id", getGame)
	api.DELETE("/games/:id", deleteGame)

	// publishers
	api.GET("/publishers", getPublishers)
	api.POST("/publishers", createPublisher)

	// genre
	api.GET("/genres", getGenres)
	api.POST("/genres", createGenre)

	// games- genre
	api.POST("/games/:id/genres", gameAddGenre)

	// platform
	api.GET("/platforms", getPlatforms)
	api.POST("/platforms", createPlatform)

	// games- platforms
	api.POST("/games/:id/platforms", gameAddPlatform)
}
