package routes

import (
	"fmt"
	"net/http"

	"github.com/game-api/models"
	"github.com/gin-gonic/gin"
)

func createPublisher(ctx *gin.Context) {

	var publisher models.Publisher

	err := ctx.ShouldBindJSON(&publisher)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse!"})
		return
	}

	err = publisher.Save()

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create", "err": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Publisher created!", "publisher": publisher})

}
