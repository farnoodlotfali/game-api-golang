package models

import (
	"database/sql"
	"fmt"
	"strconv"
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
	ImageUrl     string    `json:"image_url"`
}

func CountPublishers(q string) (int64, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM publishers WHERE title LIKE %s", "'%"+q+"%'")
	var totalCount int64
	err := db.DB.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func GetAllPublishers(page, limit, order, q, sort string) (PageResponseType[[]Publisher], error) {
	total, err := CountPublishers(q)
	if err != nil {
		return PageResponseType[[]Publisher]{}, err
	}

	fmt.Println("Total publishers count:", total)
	query := `SELECT * FROM publishers WHERE 1=1`

	emptyArr := PageResponseType[[]Publisher]{}

	// search by title
	if q != "" {
		query += fmt.Sprintf(" AND title LIKE %s", "'%"+q+"%'")
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

		full := objS3.GetS3Endpoint() + rawCover.String
		publisher.ImageUrl = full

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
	query := `INSERT INTO publishers (title, country, founding_date, website_url, image_url) VALUES ($1, $2, $3, $4, $5)
	RETURNING id`

	err := db.DB.QueryRow(
		query,
		p.Title,
		p.Country,
		p.FoundingDate,
		p.WebsiteUrl,
		p.ImageUrl,
	).Scan(&p.ID)

	return err
}
