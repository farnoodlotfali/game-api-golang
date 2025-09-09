package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/game-api/db"
)

type Genre struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

func CountGenres(query string) (int64, error) {
	countQuery := strings.Replace(query, "SELECT *", "SELECT COUNT(*)", 1)

	var totalCount int64
	err := db.DB.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func GetAllGenres(page, limit, order, q, sort string) (PageResponseType[[]Genre], error) {

	query := `SELECT * FROM genres WHERE 1=1`

	emptyArr := PageResponseType[[]Genre]{}

	// search by name
	if q != "" {
		query += fmt.Sprintf(" AND name ILIKE %s", "'%"+q+"%'")
	}

	// order and sort
	if sort != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sort, order)
	}

	total, err := CountGenres(query)
	if err != nil {
		return emptyArr, err
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
	} else if lastPage == 0 {
		lastPage = 1
	}

	rows, err := db.DB.Query(query)
	if err != nil {
		return emptyArr, err
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
			return emptyArr, err
		}
		genres = append(genres, genre)
	}
	if len(genres) == 0 {
		return SuccessPaginationResponse([]Genre{}, total, lastPage, intPage)
	}

	return SuccessPaginationResponse(genres, total, lastPage, intPage)
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

func (g *Genre) GetID() int64 {
	return g.ID
}
