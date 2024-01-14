package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/duduBTW/steam-ssr/db"
	"github.com/duduBTW/steam-ssr/presentation/game"
	"github.com/duduBTW/steam-ssr/presentation/home"
	"github.com/labstack/echo/v4"
)

type GameHandler struct {
	Database *sql.DB
}

func (h GameHandler) GamePage(c echo.Context) error {
	gameId := c.Param("gameId")
	gameData := db.GameDb{Database: h.Database}.GetGame(gameId)

	userId, _ := Utils{}.GetCookieValueAsInt("Authentication", c)
	cartId, _ := Utils{}.GetCookieValueAsInt("CartIdentifier", c)

	isGameOnCart := db.IsGameOnCart(gameId, userId, cartId, h.Database)

	return game.GamePage(game.Props{
		Game:         gameData,
		IsGameOnCart: isGameOnCart,
	}).Render(c.Request().Context(), c.Response())
}

func (h GameHandler) FastSearch(c echo.Context) error {
	gamesData := db.GameDb{Database: h.Database}.SearchGameByTitle(c.QueryParam("searchKey"))

	if len(gamesData) == 0 {
		return c.String(http.StatusNotFound, "")
	}

	return game.FastSearch(game.FastSearchProps{
		Games: gamesData,
	}).Render(c.Request().Context(), c.Response())
}

type WishlistHandler struct {
	Database *sql.DB
}

func (h WishlistHandler) WishlistButton(c echo.Context) error {
	time.Sleep(3 * time.Second)
	return home.WishlistedButton().Render(c.Request().Context(), c.Response())
}
