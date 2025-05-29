package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/game-api/models"
	"github.com/game-api/objS3"
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
func getGame(ctx *gin.Context) {
	gameId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch!"})
		return
	}

	game, err := models.GetGameByID(gameId)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch!"})
		return
	}
	ctx.JSON(http.StatusOK, game)
}

func createGame_withRowJson(ctx *gin.Context) {

	var game models.GameCreateDTO

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

	err = game.GameUpdateGenre(game.GenreIds)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add genre to game"})
		return
	}

	err = game.GameUpdatePlatform(game.PlatformsIds)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add platform to game"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Game created!", "game": game})

}

func createGame(ctx *gin.Context) {
	// 1) Parse form fields
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid form"})
		return
	}

	var dto models.GameCreateDTO
	dto.Title = ctx.Request.FormValue("title")
	dto.PublisherID, _ = strconv.ParseInt(ctx.Request.FormValue("publisher_id"), 10, 64)
	dto.ReleaseDate, _ = time.Parse(time.RFC3339, ctx.Request.FormValue("release_date"))
	dto.Description = ptrString(ctx.Request.FormValue("description"))

	// Parse array fields
	for _, s := range ctx.Request.Form["genre_ids[]"] {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			dto.GenreIds = append(dto.GenreIds, id)
		}
	}
	for _, s := range ctx.Request.Form["platform_ids[]"] {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			dto.PlatformsIds = append(dto.PlatformsIds, id)
		}
	}

	// 2) Handle cover image upload
	if fileHeader, err := ctx.FormFile("cover_image_url"); err == nil {
		url, err := objS3.UploadFileToS3(fileHeader, "test")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to upload to S3"})
			return
		}
		dto.CoverImageURL = &url
	}

	// 3) Save to database
	if err := dto.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create game", "err": err.Error()})
		return
	}

	// 4) Associate genres & platforms
	dto.GameUpdateGenre(dto.GenreIds)
	dto.GameUpdatePlatform(dto.PlatformsIds)

	ctx.JSON(http.StatusOK, gin.H{"message": "Game created!", "game": dto})
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func updateGame(ctx *gin.Context) {
	gameId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse game id"})
		return
	}
	_, err = models.GetGameByID(gameId)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch game"})
		return
	}

	var updatedGame models.GameCreateDTO

	err = ctx.ShouldBindJSON(&updatedGame)
	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse!"})
		return
	}
	updatedGame.ID = gameId

	err = updatedGame.Update()
	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update", "err": err})
		return
	}

	err = updatedGame.GameUpdateGenre(updatedGame.GenreIds)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update genre of game"})
		return
	}

	err = updatedGame.GameUpdatePlatform(updatedGame.PlatformsIds)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update platform of game"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Game updated!", "game": updatedGame})

}
