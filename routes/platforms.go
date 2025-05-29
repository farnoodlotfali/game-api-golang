package routes

import (
	"net/http"

	"github.com/game-api/models"
	"github.com/gin-gonic/gin"
)

func getPlatforms(ctx *gin.Context) {
	platforms, err := models.GetAllPlatforms()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch!"})
		return
	}
	ctx.JSON(http.StatusOK, platforms)
}

func createPlatform(ctx *gin.Context) {
	var platform models.Platform

	err := ctx.ShouldBindJSON(&platform)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse!"})
		return
	}

	err = platform.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create", "err": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "platform created!", "platform": platform})

}
