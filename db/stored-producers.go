package db

func createGameGenreProducers() {
	sql := `
	CREATE OR REPLACE PROCEDURE replace_game_genres(g_id BIGINT, g_genres BIGINT[])
	LANGUAGE plpgsql
	AS $$
	BEGIN
		DELETE FROM game_genres WHERE game_id = g_id;

		INSERT INTO game_genres (game_id, genre_id)
		SELECT g_id, unnest(g_genres)
		ON CONFLICT (game_id, genre_id) DO NOTHING;
	END;
	$$;
	`

	_, err := DB.Exec(sql)
	if err != nil {
		panic("Failed to create replace_game_genres stored procedure: ")
	}

}

func createGamePlatformProducers() {
	sql := `
	CREATE OR REPLACE PROCEDURE replace_game_platforms(g_id BIGINT, g_platforms BIGINT[])
	LANGUAGE plpgsql
	AS $$
	BEGIN
		DELETE FROM game_platforms WHERE game_id = g_id;

		INSERT INTO game_platforms (game_id, platform_id)
		SELECT g_id, unnest(g_platforms)
		ON CONFLICT (game_id, platform_id) DO NOTHING;
	END;
	$$;
	`

	_, err := DB.Exec(sql)
	if err != nil {
		panic("Failed to create replace_game_platforms stored procedure: ")
	}

}
