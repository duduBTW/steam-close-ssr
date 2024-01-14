package db

import (
	"database/sql"
	"log"

	"github.com/duduBTW/steam-ssr/data"
)

type Game struct {
	CoverImage  string
	Images      []string
	Name        string
	Description string
	Price       string
	GameId      int64
}

type GameDb struct {
	Database *sql.DB
}

func (gameDb GameDb) GetGame(gameId string) data.Game {
	query := `SELECT GameId, Title, Description, CoverImage, Price FROM Game WHERE GameId = $1`
	var game data.Game

	err := gameDb.Database.QueryRow(query, gameId).Scan(&game.GameId, &game.Name, &game.Description, &game.CoverImage, &game.Price)
	if err != nil {
		log.Fatal(err)
	}

	game.Images = gameDb.GetGameCovers(gameId)
	return game
}

func (gameDb GameDb) GetGameCovers(gameId string) []data.Image {
	query := `
		SELECT
			Image.ImageId,
			Image.Medium,
			Image.Small,
			Image.Large
		FROM
			Image
		JOIN
			GameImages ON Image.ImageId = GameImages.ImageId
		JOIN
			Game ON Game.GameId = GameImages.GameId
		WHERE
			Game.GameId = $1
	`
	var covers []data.Image

	rows, err := gameDb.Database.Query(query, gameId)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var cover data.Image

		err := rows.Scan(&cover.ImageId, &cover.Medium, &cover.Small, &cover.Large)
		if err != nil {
			log.Fatal(err)
		}

		covers = append(covers, cover)
	}

	return covers
}

func (gameDb GameDb) SearchGameByTitle(searchKey string) []data.Game {
	query := `
		SELECT 
			GameId, 
			Title, 
			Description, 
			CoverImage, 
			Price
		FROM 
			Game 
		WHERE 
			Game.Title LIKE '%' || $1 ||'%'
	`
	var games []data.Game

	rows, err := gameDb.Database.Query(query, searchKey)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var game data.Game

		err := rows.Scan(&game.GameId, &game.Name, &game.Description, &game.CoverImage, &game.Price)
		if err != nil {
			log.Fatal(err)
		}

		games = append(games, game)
	}

	return games
}
