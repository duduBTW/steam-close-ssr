package handlers

import (
	"database/sql"
	"log"

	"github.com/duduBTW/steam-ssr/data"
	"github.com/duduBTW/steam-ssr/presentation/home"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
	Database *sql.DB
}

func (handler HomeHandler) HomePage(context echo.Context) error {
	return home.HomePage(home.Props{
		Games: getList(handler.Database),
	}).Render(context.Request().Context(), context.Response())
}

func getList(db *sql.DB) []data.Game {
	var games []data.Game
	gamesToIdMap := make(map[int]data.Game)

	query := `
		SELECT 
			Game.GameId, 
			Game.Title, 
			Game.Description, 
			Game.CoverImage, 
			Game.Price,
			Image.ImageId,
			Image.Medium,
			Image.Small,
			Image.Large
		FROM 
			Game
		JOIN (
			SELECT 
					GameId,
					ImageId,
					ROW_NUMBER() OVER (PARTITION BY GameId ORDER BY ImageId) AS RowNum
			FROM 
					GameImages
		) AS GameImages ON Game.GameId = GameImages.GameId
		JOIN
			Image ON Image.ImageId = GameImages.ImageId
		WHERE
			GameImages.RowNum <= 4;
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var gameRow data.Game
		var image data.Image

		err := rows.Scan(
			&gameRow.GameId,
			&gameRow.Name,
			&gameRow.Description,
			&gameRow.CoverImage,
			&gameRow.Price,
			&image.ImageId,
			&image.Medium,
			&image.Small,
			&image.Large,
		)
		if err != nil {
			log.Fatal(err)
		}

		_, hasGame := gamesToIdMap[gameRow.GameId]
		if !hasGame {
			games = append(games, gameRow)
			gamesToIdMap[gameRow.GameId] = gameRow
		}

		games[len(games)-1].Images = append(games[len(games)-1].Images, image)
	}

	return games
}
