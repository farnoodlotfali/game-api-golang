package models

import "github.com/game-api/db"

type Platform struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

func GetAllPlatforms() ([]Platform, error) {
	query := `SELECT * FROM platforms`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var platforms []Platform
	for rows.Next() {
		var platform Platform
		err := rows.Scan(
			&platform.ID,
			&platform.Name,
			&platform.Description,
		)
		if err != nil {
			return nil, err
		}
		platforms = append(platforms, platform)
	}
	return platforms, nil
}

func (p *Platform) Save() error {
	query := `INSERT INTO platforms (name, description) VALUES ($1, $2) RETURNING id`

	err := db.DB.QueryRow(
		query,
		p.Name,
		p.Description,
	).Scan(&p.ID)

	return err
}
