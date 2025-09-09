package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func badRequest400(ctx *gin.Context, text string, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{"message": text, "err": err.Error()})

}

func notFoundRequest404(ctx *gin.Context, text string, err error) {
	ctx.JSON(http.StatusNotFound, gin.H{"message": text, "err": err.Error()})

}

func serverErrorRequest500(ctx *gin.Context, text string, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"message": text, "err": err.Error()})

}
