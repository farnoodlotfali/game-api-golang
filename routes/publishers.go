package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/game-api/models"
	"github.com/game-api/objS3"
	"github.com/gin-gonic/gin"
)

func getPublishers(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")     // Default to 1 if not provided
	limit := ctx.DefaultQuery("limit", "10")  // Default to 10 if not provided
	q := ctx.DefaultQuery("q", "")            // Default sorting by title
	order := ctx.DefaultQuery("order", "asc") // Default sorting order is ascending
	sort := ctx.DefaultQuery("sort", "")      // Default sorting order is ascending

	games, err := models.GetAllPublishers(page, limit, order, q, sort)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch!", "err": err})
		return
	}
	ctx.JSON(http.StatusOK, games)
}

func createPublisher(ctx *gin.Context) {
	// 1) Parse form fields
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid form", "err": err.Error()})
		return
	}

	var publisher models.Publisher
	publisher.Title = ctx.Request.FormValue("title")
	publisher.Country = ctx.Request.FormValue("country")

	// Handle optional website_url
	websiteURL := ctx.Request.FormValue("website_url")
	if websiteURL != "" {
		publisher.WebsiteUrl = websiteURL
	}

	date, err := time.Parse(time.RFC3339, ctx.Request.FormValue("founding_date"))

	if err == nil {
		publisher.FoundingDate = date
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "date error", "err": err.Error()})
		return
	}

	// 2) Handle  image upload
	fileHeader, err := ctx.FormFile("image_url")
	if err == nil {
		// If a file is uploaded, upload it to S3
		url, err := objS3.UploadFileToS3(fileHeader, "test", "publishers/")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to upload to S3", "err": err.Error()})
			return
		}
		publisher.ImageUrl = url
	}

	// 3) Save to database
	if err := publisher.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create publisher", "err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Publisher created!", "publisher": publisher})

}
