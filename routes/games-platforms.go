package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/game-api/models"
	"github.com/gin-gonic/gin"
)

type AddPlatformRequest struct {
	PlatformID []int64 `json:"platform_id" binding:"required"`
}

func gameAddPlatform(ctx *gin.Context) {
	gameId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse game id"})
		return
	}

	var req AddPlatformRequest
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

	err = game.GameUpdatePlatform(req.PlatformID)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add platform to game"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "added!"})

}
