package routes

import (
	"net/http"

	"github.com/game-api/models"
	"github.com/gin-gonic/gin"
)

func getPlatforms(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")     // Default to 1 if not provided
	limit := ctx.DefaultQuery("limit", "10")  // Default to 10 if not provided
	q := ctx.DefaultQuery("q", "")            // Default sorting by title
	order := ctx.DefaultQuery("order", "asc") // Default sorting order is ascending
	sort := ctx.DefaultQuery("sort", "")      //

	platforms, err := models.GetAllPlatforms(page, limit, order, q, sort)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch!", "err": err.Error()})
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
