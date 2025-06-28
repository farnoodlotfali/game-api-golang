package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/game-api/db"
	"github.com/game-api/objS3"
)

type Game struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title" binding:"required"`
	ReleaseDate   time.Time `json:"release_date" binding:"required"`
	CoverImageURL *string   `json:"cover_image_url"`
	Description   *string   `json:"description"`
	PublisherID   int64     `json:"publisher_id" binding:"required"`
}
type GameCreateDTO struct {
	Game
	GenreIds      []int64  `json:"genre_ids"  binding:"required"`
	PlatformsIds  []int64  `json:"platform_ids"  binding:"required"`
	ScreenshotIds *[]int64 `json:"screenshot_ids"`
}
type GameDTO struct {
	Game
	Publisher   Publisher     `json:"publisher"`
	Genres      []Genre       `json:"genres"`
	Platforms   []Platform    `json:"platforms"`
	Screenshots *[]Screenshot `json:"screenshots"`
}

func CountGames(q string) (int64, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM game_full_info WHERE title LIKE %s", "'%"+q+"%'")
	var totalCount int64
	err := db.DB.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func GetAllGames(page, limit, order, q, sort string) (PageResponseType[[]GameDTO], error) {

	total, err := CountGames(q)
	if err != nil {
		return PageResponseType[[]GameDTO]{}, err
	}

	// fmt.Println("Total games count:", total)
	query := `SELECT * FROM game_full_info WHERE 1=1`

	emptyArr := PageResponseType[[]GameDTO]{}

	// search by title
	if q != "" {
		query += fmt.Sprintf(" AND title ILIKE %s", "'%"+q+"%'")
	}

	// order and sort
	if sort != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sort, order)
	}

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)

	offset := (intPage - 1) * intLimit

	// limit
	query += fmt.Sprintf(" LIMIT %s", limit)

	// offset
	query += fmt.Sprintf(" OFFSET %s", fmt.Sprint(offset))

	// lastPage
	lastPage := int(total) / intLimit
	if int(total)%intLimit != 0 {
		lastPage++
	} else {
		lastPage = 1
	}

	// fmt.Println(query)

	rows, err := db.DB.Query(query)
	if err != nil {
		return emptyArr, err
	}
	defer rows.Close()

	var games []GameDTO
	var rawCover sql.NullString
	for rows.Next() {
		var (
			game            GameDTO
			publisherJSON   []byte
			genresJSON      []byte
			platformsJSON   []byte
			screenshotsJSON []byte
		)

		err := rows.Scan(
			&game.ID,
			&game.Title,
			&game.ReleaseDate,
			&rawCover,
			&game.Description,
			&game.PublisherID,
			&publisherJSON,
			&genresJSON,
			&platformsJSON,
			&screenshotsJSON,
		)

		full := objS3.GetS3Endpoint() + rawCover.String
		game.CoverImageURL = &full

		if err != nil {
			return emptyArr, err
		}

		if err := json.Unmarshal(publisherJSON, &game.Publisher); err != nil {
			return emptyArr, err
		}

		if len(genresJSON) == 0 || string(genresJSON) == "null" {
			game.Genres = []Genre{}
		} else if err := json.Unmarshal(genresJSON, &game.Genres); err != nil {
			return emptyArr, err
		}

		if len(platformsJSON) == 0 || string(platformsJSON) == "null" {
			game.Platforms = []Platform{}
		} else if err := json.Unmarshal(platformsJSON, &game.Platforms); err != nil {
			return emptyArr, err
		}

		if len(screenshotsJSON) == 0 || string(screenshotsJSON) == "null" {
			game.Screenshots = &[]Screenshot{}
		} else if err := json.Unmarshal(screenshotsJSON, &game.Screenshots); err != nil {
			return emptyArr, err
		}

		games = append(games, game)
	}

	if len(games) == 0 {
		return SuccessPaginationResponse([]GameDTO{}, total, lastPage, intPage)
	}

	return SuccessPaginationResponse(games, total, lastPage, intPage)

}

func GetGameByID(id int64) (*GameDTO, error) {
	query := `SELECT * FROM game_full_info WHERE game_id = $1`

	row := db.DB.QueryRow(query, id)
	var (
		game            GameDTO
		publisherJSON   []byte
		genresJSON      []byte
		platformsJSON   []byte
		screenshotsJSON []byte
		rawCover        sql.NullString
	)

	err := row.Scan(
		&game.ID,
		&game.Title,
		&game.ReleaseDate,
		&rawCover,
		&game.Description,
		&game.PublisherID,
		&publisherJSON,
		&genresJSON,
		&platformsJSON,
		&screenshotsJSON,
	)

	full := objS3.GetS3Endpoint() + rawCover.String
	game.CoverImageURL = &full

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(publisherJSON, &game.Publisher); err != nil {
		return nil, err
	}

	if len(genresJSON) == 0 || string(genresJSON) == "null" {
		game.Genres = []Genre{}
	} else if err := json.Unmarshal(genresJSON, &game.Genres); err != nil {
		return nil, err
	}

	if len(platformsJSON) == 0 || string(platformsJSON) == "null" {
		game.Platforms = []Platform{}
	} else if err := json.Unmarshal(platformsJSON, &game.Platforms); err != nil {
		return nil, err
	}

	if len(screenshotsJSON) == 0 || string(screenshotsJSON) == "null" {
		game.Screenshots = &[]Screenshot{}
	} else if err := json.Unmarshal(screenshotsJSON, &game.Screenshots); err != nil {
		return nil, err
	}

	return &game, nil

}

func (g *Game) Save() error {

	query := `INSERT INTO games (title, release_date, cover_image_url, description, publisher_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := db.DB.QueryRow(
		query,
		g.Title,
		g.ReleaseDate,
		g.CoverImageURL,
		g.Description,
		g.PublisherID,
	).Scan(&g.ID)

	return err

}

func (g *Game) Update() error {

	query := `
	UPDATE events 
    SET 
        title = $1,
        release_date = $2,
        cover_image_url = $3,
        description = $4
        publisher_id = $5
    WHERE 
        id = $6
	`

	_, err := db.DB.Exec(
		query,
		g.Title,
		g.ReleaseDate,
		g.CoverImageURL,
		g.Description,
		g.PublisherID,
		g.ID,
	)

	return err

}

func (g *Game) Delete() error {

	query := "CALL delete_game($1)"

	_, err := db.DB.Exec(query, g.ID)
	return err

}
