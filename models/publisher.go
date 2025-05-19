package models

import (
	"time"

	"github.com/game-api/db"
)

type Publisher struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title" binding:"required"`
	Country      string    `json:"country" binding:"required"`
	FoundingDate time.Time `json:"founding_date" binding:"required"`
	WebsiteUrl   string    `json:"website_url"`
	ImageUrl     string    `json:"image_url"`
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
