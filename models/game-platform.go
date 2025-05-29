package models

import (
	"github.com/game-api/db"
	"github.com/lib/pq"
)

func (g *Game) GameUpdatePlatform(platformIDs []int64) error {

	query := "CALL replace_game_platforms($1, $2)"

	_, err := db.DB.Exec(query, g.ID, pq.Array(platformIDs))
	return err
}
