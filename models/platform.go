package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/game-api/db"
)

type Platform struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

func CountPlatforms(query string) (int64, error) {
	countQuery := strings.Replace(query, "SELECT *", "SELECT COUNT(*)", 1)

	var totalCount int64
	err := db.DB.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func GetAllPlatforms(page, limit, order, q, sort string) (PageResponseType[[]Platform], error) {
	emptyArr := PageResponseType[[]Platform]{}

	query := `SELECT * FROM platforms WHERE 1=1`

	// search by name
	if q != "" {
		query += fmt.Sprintf(" AND name ILIKE %s", "'%"+q+"%'")
	}

	// order and sort
	if sort != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sort, order)
	}

	total, err := CountPlatforms(query)
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

	fmt.Println(query)

	rows, err := db.DB.Query(query)
	if err != nil {
		return emptyArr, err
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
			return emptyArr, err
		}
		platforms = append(platforms, platform)
	}

	if len(platforms) == 0 {
		return SuccessPaginationResponse([]Platform{}, total, lastPage, intPage)
	}

	return SuccessPaginationResponse(platforms, total, lastPage, intPage)
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
func (g *Platform) GetID() int64 {
	return g.ID
}
