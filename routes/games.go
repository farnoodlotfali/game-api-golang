package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/game-api/models"
	"github.com/game-api/objS3"
	"github.com/game-api/utils"
	"github.com/gin-gonic/gin"
)

func getGames(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")
	q := ctx.DefaultQuery("q", "")
	order := ctx.DefaultQuery("order", "asc")
	sort := ctx.DefaultQuery("sort", "")
	releaseDateFrom := ctx.DefaultQuery("releaseDateFrom", "")
	releaseDateTo := ctx.DefaultQuery("releaseDateTo", "")

	genres_ids, err := utils.ExtractIDsFromQuery[*models.Genre](ctx, "genre_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	platform_ids, err := utils.ExtractIDsFromQuery[*models.Platform](ctx, "platform_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	publisher_ids, err := utils.ExtractIDsFromQuery[*models.Publisher](ctx, "publisher_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	games, err := models.GetAllGames(page, limit, order, q, sort, releaseDateFrom, releaseDateTo, genres_ids, platform_ids, publisher_ids)
	if err != nil {
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
	if *game.CoverImage != "" {
		full := objS3.AddEndPointToUrl(*game.CoverImage)
		game.CoverImage = &full
	}

	if err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Could not Found!"})
		return
	}

	gameData, _ := models.SuccessResponse(game)
	ctx.JSON(http.StatusOK, gameData)
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
	dto.Description = utils.PtrString(ctx.Request.FormValue("description"))

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
	fileHeader, err := ctx.FormFile("cover_image")

	if err == nil {

		url, err := objS3.UploadFileToS3(fileHeader, "test", "games/")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to upload to S3", "err": err.Error()})
			return
		}
		dto.CoverImage = &url
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

func updateGame(ctx *gin.Context) {
	// 1) Parse form fields
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		badRequest400(ctx, "invalid form", err)
		return
	}

	gameId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		badRequest400(ctx, "Could not parse game id", err)
		return
	}
	foundGame, err := models.GetGameByID(gameId)

	if err != nil {
		notFoundRequest404(ctx, "Could not fetch game", err)
		return
	}

	var updatedGame models.GameCreateDTO
	updatedGame.Title = ctx.Request.FormValue("title")
	updatedGame.PublisherID, _ = strconv.ParseInt(ctx.Request.FormValue("publisher_id"), 10, 64)
	updatedGame.ReleaseDate, _ = time.Parse(time.RFC3339, ctx.Request.FormValue("release_date"))
	updatedGame.Description = utils.PtrString(ctx.Request.FormValue("description"))
	updatedGame.ID = gameId
	updatedGame.CoverImage = foundGame.CoverImage

	genreIdsStr := ctx.PostFormArray("genre_ids[]")
	var genreIds []int64
	for _, s := range genreIdsStr {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			badRequest400(ctx, "Invalid genre_id: "+s, err)
			return
		}
		genreIds = append(genreIds, id)
	}
	updatedGame.GenreIds = genreIds

	platformIdsStr := ctx.PostFormArray("platform_ids[]")
	var platformIds []int64
	for _, s := range platformIdsStr {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			badRequest400(ctx, "Invalid platform_ids: "+s, err)
			return
		}
		platformIds = append(platformIds, id)
	}
	updatedGame.PlatformsIds = platformIds

	fileHeader, err := ctx.FormFile("cover_image")
	if err == nil {
		err = objS3.DeleteFileFromS3(*updatedGame.CoverImage)
		if err != nil {
			serverErrorRequest500(ctx, "Could not delete cover image", err)
			return
		}

		url, err := objS3.UploadFileToS3(fileHeader, "test", "games/")
		if err != nil {
			serverErrorRequest500(ctx, "failed to upload to S3", err)
			return
		}
		updatedGame.CoverImage = &url
	}

	err = updatedGame.Update()
	if err != nil {
		serverErrorRequest500(ctx, "Could not update", err)
		return
	}

	err = updatedGame.GameUpdateGenre(updatedGame.GenreIds)
	if err != nil {
		serverErrorRequest500(ctx, "Could not update genre of game", err)
		return
	}

	err = updatedGame.GameUpdatePlatform(updatedGame.PlatformsIds)
	if err != nil {
		serverErrorRequest500(ctx, "Could not update platform of game", err)
		return
	}

	if *updatedGame.CoverImage != "" {
		full := objS3.AddEndPointToUrl(*updatedGame.CoverImage)
		updatedGame.CoverImage = &full
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
