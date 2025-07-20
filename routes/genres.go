package routes

import (
	"fmt"
	"net/http"

	"github.com/game-api/models"
	"github.com/gin-gonic/gin"
)

func getGenres(ctx *gin.Context) {

	page := ctx.DefaultQuery("page", "1")     // Default to 1 if not provided
	limit := ctx.DefaultQuery("limit", "10")  // Default to 10 if not provided
	q := ctx.DefaultQuery("q", "")            // Default sorting by title
	order := ctx.DefaultQuery("order", "asc") // Default sorting order is ascending
	sort := ctx.DefaultQuery("sort", "")      //

	genres, err := models.GetAllGenres(page, limit, order, q, sort)

	if err != nil {
		fmt.Print(err)
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
