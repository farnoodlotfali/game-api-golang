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

	page := ctx.DefaultQuery("page", "1")     // Default to 1 if not provided
	limit := ctx.DefaultQuery("limit", "10")  // Default to 10 if not provided
	q := ctx.DefaultQuery("q", "")            // Default sorting by title
	order := ctx.DefaultQuery("order", "asc") // Default sorting order is ascending
	sort := ctx.DefaultQuery("sort", "")      // Default sorting order is ascending

	fmt.Println("Page:", page, "Limit:", limit, "q:", q, "Order:", order, "sort:", sort)

	games, err := models.GetAllGames(page, limit, order, q, sort)

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
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Could not Found!"})
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
	fileHeader, err := ctx.FormFile("cover_image_url")

	if err == nil {

		url, err := objS3.UploadFileToS3(fileHeader, "test", "games/")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to upload to S3", "err": err.Error()})
			return
		}
		dto.CoverImageURL = &url
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "UploadFileToS3 error", "err": err.Error()})
		return
	}

	// 3) Save to database
	if err := dto.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create game", "err": err.Error()})
		return
	}

	// 4) Associate genres & platforms
	dto.GameUpdateGenre(dto.GenreIds)
	dto.GameUpdatePlatform(dto.PlatformsIds)

	screenshotIds := make([]int64, 0)
	form := ctx.Request.MultipartForm
	if form != nil && form.File != nil {
		if screenshotHeaders, exists := form.File["screenshot[]"]; exists {
			for _, sh := range screenshotHeaders {
				s3url, err := objS3.UploadFileToS3(sh, "test", "screenshots/")
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"message": "failed to upload screenshot to S3",
						"error":   err.Error(),
					})
					return
				}

				shot := models.Screenshot{
					Url:    s3url,
					GameID: dto.ID,
				}
				if err := shot.Save(); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"message": "could not save screenshot record",
						"error":   err.Error(),
					})
					return
				}
				screenshotIds = append(screenshotIds, shot.ID)

			}
		}
	}

	dto.ScreenshotIds = &screenshotIds

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

func deleteGame(ctx *gin.Context) {
	gameId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse game id"})
		return
	}
	var game *models.GameDTO
	game, err = models.GetGameByID(gameId)

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch game"})
		return
	}

	err = game.Delete()

	for _, screenshot := range *game.Screenshots {
		err = objS3.DeleteFileFromS3(screenshot.Url)

		if err != nil {
			fmt.Print(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete screenshot"})
			return
		}
	}

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete game"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Game deleted!", "game": game})

}
