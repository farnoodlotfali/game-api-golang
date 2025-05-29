package db

import "fmt"

func createGameView() {
	sql := `
	DROP VIEW IF EXISTS game_full_info;
	CREATE OR REPLACE VIEW game_full_info AS
	SELECT
		g.id AS game_id,
		g.title,
		g.release_date,
		g.cover_image_url,
		g.description,
		g.publisher_id,

		JSONB_BUILD_OBJECT(
			'id', pub.id,
			'title', pub.title,
			'country', pub.country,
			'founding_date', TO_CHAR(pub.founding_date AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			'website_url', pub.website_url,
			'image_url', pub.image_url
		) AS publisher,

		JSON_AGG(DISTINCT JSONB_BUILD_OBJECT(
			'id', ge.id,
			'name', ge.name,
			'description', ge.description
		)) FILTER (WHERE ge.id IS NOT NULL) AS genres,

		JSON_AGG(DISTINCT JSONB_BUILD_OBJECT(
			'id', pl.id,
			'name', pl.name,
			'description', pl.description
		)) FILTER (WHERE pl.id IS NOT NULL) AS platforms,

		JSON_AGG(DISTINCT JSONB_BUILD_OBJECT(
			'id', sc.id,
			'game_id', sc.game_id,
			'url', sc.url
		)) FILTER (WHERE sc.id IS NOT NULL) AS screenshots

	FROM games g
	JOIN publishers pub ON g.publisher_id = pub.id
	LEFT JOIN screenshots sc ON g.id = sc.game_id
	LEFT JOIN game_genres gg ON g.id = gg.game_id
	LEFT JOIN genres ge ON gg.genre_id = ge.id
	LEFT JOIN game_platforms gp ON g.id = gp.game_id
	LEFT JOIN platforms pl ON gp.platform_id = pl.id
	GROUP BY g.id, pub.id;

	`

	_, err := DB.Exec(sql)
	if err != nil {
		fmt.Print(err)
		panic("Failed to create game_full_info view: ")
	}

}
