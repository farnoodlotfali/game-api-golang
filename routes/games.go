package routes

import (
	"fmt"
	"net/http"

	"github.com/game-api/models"
	"github.com/gin-gonic/gin"
)

func getGames(ctx *gin.Context) {
	games, err := models.GetAllGames()

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch!"})
		return
	}
	ctx.JSON(http.StatusOK, games)
}

func createGame(ctx *gin.Context) {

	var game models.Game

	err := ctx.ShouldBindJSON(&game)
	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse!"})
		return
	}

	err = game.Save()
	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create", "err": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Game created!", "game": game})

}
