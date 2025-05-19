package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "game"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		panic("Error connecting to the database ")
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	createTables()
}
func createTables() {

	createGamesTable()
	createPublishersTable()

}

func createGamesTable() {
	query := `CREATE TABLE IF NOT EXISTS games (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		release_date TIMESTAMP NOT NULL,
		cover_image_url VARCHAR(512),
		description VARCHAR(512),
		publisher_id INTEGER,
		CONSTRAINT fk_user
            FOREIGN KEY(publisher_id) 
            REFERENCES publishers(id)
            ON DELETE SET NULL
	)`

	_, err := DB.Exec(query)

	if err != nil {
		panic("cannot create games table" + err.Error())
	}

}

func createPublishersTable() {
	query := `CREATE TABLE IF NOT EXISTS publishers (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		country VARCHAR(30) NOT NULL,
		founding_date TIMESTAMP NOT NULL,
		website_url VARCHAR(255),
		image_url VARCHAR(512)
	)
	`
	_, err := DB.Exec(query)

	if err != nil {
		panic("cannot create publishers table" + err.Error())
	}

}
