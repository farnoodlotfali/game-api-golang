package models

import "github.com/game-api/db"

type Genre struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

func GetAllGenres() ([]Genre, error) {
	query := `SELECT * FROM genres`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []Genre
	for rows.Next() {
		var genre Genre
		err := rows.Scan(
			&genre.ID,
			&genre.Name,
			&genre.Description,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	if len(genres) == 0 {
		return []Genre{}, nil
	}

	return genres, nil
}

func (g *Genre) Save() error {

	query := `INSERT INTO genres (name, description) VALUES ($1, $2) RETURNING id`

	err := db.DB.QueryRow(
		query,
		g.Name,
		g.Description,
	).Scan(&g.ID)

	return err
}
