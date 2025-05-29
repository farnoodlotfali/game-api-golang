package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/game-api/models"
	"github.com/gin-gonic/gin"
)

type AddGenreRequest struct {
	GenreID []int64 `json:"genre_id" binding:"required"`
}

func gameAddGenre(ctx *gin.Context) {
	gameId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse game id"})
		return
	}

	var req AddGenreRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	game, err := models.GetGameByID(gameId)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch game"})
		return
	}

	err = game.GameUpdateGenre(req.GenreID)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add genre to game"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "added!"})

}
