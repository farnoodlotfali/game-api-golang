package routes

import (
	"net/http"

	"github.com/game-api/models"
	"github.com/gin-gonic/gin"
)

func getGenres(ctx *gin.Context) {
	genres, err := models.GetAllGenres()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch!"})
		return
	}
	ctx.JSON(http.StatusOK, genres)
}

func createGenre(ctx *gin.Context) {
	var genre models.Genre

	err := ctx.ShouldBindJSON(&genre)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse!"})
		return
	}

	err = genre.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create", "err": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Genre created!", "genre": genre})

}
