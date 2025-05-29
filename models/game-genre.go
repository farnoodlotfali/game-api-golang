package models

import (
	"github.com/game-api/db"
	"github.com/lib/pq"
)

func (g *Game) GameUpdateGenre(genreIDs []int64) error {

	query := "CALL replace_game_genres($1, $2)"

	_, err := db.DB.Exec(query, g.ID, pq.Array(genreIDs))
	return err
}
