package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/game-api/db"
	"github.com/game-api/objS3"
)

type Publisher struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title" binding:"required"`
	Country      string    `json:"country" binding:"required"`
	FoundingDate time.Time `json:"founding_date" binding:"required"`
	WebsiteUrl   string    `json:"website_url"`
	Image        string    `json:"image"`
}

func CountPublishers(query string) (int64, error) {
	countQuery := strings.Replace(query, "SELECT *", "SELECT COUNT(*)", 1)

	var totalCount int64
	err := db.DB.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func GetAllPublishers(page, limit, order, q, sort string) (PageResponseType[[]Publisher], error) {

	query := `SELECT * FROM publishers WHERE 1=1`

	emptyArr := PageResponseType[[]Publisher]{}

	// search by title
	if q != "" {
		query += fmt.Sprintf(" AND title ILIKE %s", "'%"+q+"%'")
	}

	total, err := CountPublishers(query)
	if err != nil {
		return emptyArr, err
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
	}

	rows, err := db.DB.Query(query)
	if err != nil {
		return emptyArr, err
	}
	defer rows.Close()

	var publishers []Publisher
	var publisher Publisher
	var rawCover sql.NullString
	for rows.Next() {

		err := rows.Scan(
			&publisher.ID,
			&publisher.Title,
			&publisher.Country,
			&publisher.FoundingDate,
			&publisher.WebsiteUrl,
			&rawCover,
		)

		if rawCover.String != "" {
			full := objS3.GetS3Endpoint() + rawCover.String
			publisher.Image = full
		} else {
			publisher.Image = ""
		}

		if err != nil {
			return emptyArr, err
		}

		publishers = append(publishers, publisher)
	}

	if len(publishers) == 0 {
		return SuccessPaginationResponse([]Publisher{}, total, lastPage, intPage)
	}

	return SuccessPaginationResponse(publishers, total, lastPage, intPage)

}

func (p *Publisher) Save() error {
	query := `INSERT INTO publishers (title, country, founding_date, website_url, image) VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	err := db.DB.QueryRow(
		query,
		p.Title,
		p.Country,
		p.FoundingDate,
		p.WebsiteUrl,
		p.Image,
	).Scan(&p.ID)

	return err
}
func (g *Publisher) GetID() int64 {
	return g.ID
}
