package models

import "github.com/game-api/db"

type Screenshot struct {
	ID     int64  `json:"id"`
	Url    string `json:"url" binding:"required"`
	GameID int64  `json:"game_id" binding:"required"`
}

func (s *Screenshot) Save() error {

	query := `INSERT INTO screenshots (url, game_id) VALUES ($1, $2) RETURNING id`

	err := db.DB.QueryRow(
		query,
		s.Url,
		s.GameID,
	).Scan(&s.ID)

	return err

}
