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

	createPublishersTable()
	createGenresTable()
	createPlatformsTable()
	createGamesTable()
	createGamesGenres()
	createGamesPlatforms()
	createScreenshotTable()

	// producers
	createGameGenreProducers()
	createGamePlatformProducers()
	deleteGameProducers()

	// views
	createGameView()

}

// games
func createGamesTable() {
	query := `CREATE TABLE IF NOT EXISTS games (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		release_date TIMESTAMP NOT NULL,
		cover_image_url VARCHAR(512),
		description VARCHAR(512),
		publisher_id INTEGER,
		CONSTRAINT fk_publisher
            FOREIGN KEY(publisher_id) 
            REFERENCES publishers(id)
            ON DELETE SET NULL
	)`

	_, err := DB.Exec(query)

	if err != nil {
		panic("cannot create games table" + err.Error())
	}

}

// publishers
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

// genres
func createGenresTable() {
	query := `CREATE TABLE IF NOT EXISTS genres (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description VARCHAR(512)
	)
	`
	_, err := DB.Exec(query)

	if err != nil {
		panic("cannot create genres table" + err.Error())
	}

}

// game_genres
func createGamesGenres() {
	query := `CREATE TABLE IF NOT EXISTS game_genres (
		id SERIAL PRIMARY KEY,
		game_id INTEGER,
		genre_id INTEGER,
		CONSTRAINT fk_game
            FOREIGN KEY(game_id) 
            REFERENCES games(id)
            ON DELETE SET NULL,
		CONSTRAINT fk_genre
            FOREIGN KEY(genre_id) 
            REFERENCES genres(id)
            ON DELETE SET NULL,
    	UNIQUE(game_id, genre_id)
	)`

	_, err := DB.Exec(query)

	if err != nil {
		panic("cannot create game_genres table" + err.Error())
	}

}

// platforms
func createPlatformsTable() {
	query := `CREATE TABLE IF NOT EXISTS platforms (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description VARCHAR(512)
	)
	`
	_, err := DB.Exec(query)

	if err != nil {
		panic("cannot create platforms table" + err.Error())
	}

}

// game_platforms
func createGamesPlatforms() {
	query := `CREATE TABLE IF NOT EXISTS game_platforms (
		id SERIAL PRIMARY KEY,
		game_id INTEGER,
		platform_id INTEGER,
		CONSTRAINT fk_game
            FOREIGN KEY(game_id) 
            REFERENCES games(id)
            ON DELETE SET NULL,
		CONSTRAINT fk_platform
            FOREIGN KEY(platform_id) 
            REFERENCES platforms(id)
            ON DELETE SET NULL,
    	UNIQUE(game_id, platform_id)
	)`

	_, err := DB.Exec(query)

	if err != nil {
		panic("cannot create game_platforms table" + err.Error())
	}

}

// screenshots
func createScreenshotTable() {
	query := `CREATE TABLE IF NOT EXISTS screenshots (
		id SERIAL PRIMARY KEY,
		game_id INTEGER,
		url VARCHAR(512),
		CONSTRAINT fk_game
            FOREIGN KEY(game_id) 
            REFERENCES games(id)
            ON DELETE SET NULL
	)
	`
	_, err := DB.Exec(query)

	if err != nil {
		panic("cannot create screenshots table" + err.Error())
	}

}
