package models

import (
	"time"

	"github.com/game-api/db"
)

type Game struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title" binding:"required"`
	ReleaseDate   time.Time `json:"release_date" binding:"required"`
	CoverImageURL *string   `json:"cover_image_url"`
	Description   *string   `json:"description"`
	PublisherID   int64     `json:"publisher_id" binding:"required"`
}

type GameDTO struct {
	Game
	Publisher Publisher `json:"publisher"`
}

func GetAllGames() ([]GameDTO, error) {
	// query := "SELECT id, title, description, release_date, cover_image_url, publisher_id FROM games"
	// query := "SELECT * FROM games"
	query := `
	SELECT * FROM games g
	JOIN publishers p ON g.publisher_id = p.id
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []GameDTO
	for rows.Next() {
		var game GameDTO
		err := rows.Scan(
			&game.ID,
			&game.Title,
			&game.ReleaseDate,
			&game.CoverImageURL,
			&game.Description,
			&game.PublisherID,
			&game.Publisher.ID,
			&game.Publisher.Title,
			&game.Publisher.Country,
			&game.Publisher.FoundingDate,
			&game.Publisher.ImageUrl,
			&game.Publisher.WebsiteUrl,
		)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
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
