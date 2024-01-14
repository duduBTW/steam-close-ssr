package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/duduBTW/steam-ssr/data"
	"github.com/duduBTW/steam-ssr/db"
	"github.com/duduBTW/steam-ssr/presentation/cart"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	Database *sql.DB
}

func (handler CartHandler) Page(context echo.Context) error {
	var cartList []data.CartGame

	cartId, cartIdOk := Utils{}.GetCookieValueAsInt("CartIdentifier", context)
	userId, authOk := Utils{}.GetCookieValueAsInt("Authentication", context)

	if cartIdOk || authOk {
		cartList = db.GetCartWithGamesById(cartId, userId, handler.Database)
	}

	spew.Dump(cartList)

	return cart.CartPage(cart.Props{
		Games: cartList,
		Total: "$399,99",
	}).Render(context.Request().Context(), context.Response())
}

func (handler CartHandler) InsertCart(context echo.Context) error {
	userId, authOk := Utils{}.GetCookieValueAsInt("Authentication", context)
	cartId, cartIdOk := Utils{}.GetCookieValueAsInt("CartIdentifier", context)
	gameId, err := strconv.Atoi(context.QueryParam("GameId"))
	if err != nil {
		log.Fatal(err)
	}

	var userCartId int
	if authOk {
		userCartId = db.UserCartId(userId, handler.Database)
	}

	fmt.Println(cartId, gameId)

	if (!cartIdOk && !authOk) || (authOk && userCartId == 0) {
		cartId = db.CreateCart(userId, handler.Database)
	} else if authOk {
		cartId = userCartId
	}

	db.InsertGameToCart(gameId, cartId, handler.Database)

	if !cartIdOk && !authOk {
		cartCookie := Utils{}.CreteCookie("CartIdentifier", strconv.Itoa(cartId))
		context.SetCookie(cartCookie)
	}

	return context.Redirect(http.StatusMovedPermanently, "/cart")
}

func (handler CartHandler) RemoveCartGame(context echo.Context) error {
	time.Sleep(3 * time.Second)

	userId, authOk := Utils{}.GetCookieValueAsInt("Authentication", context)
	cartGameId, err := strconv.Atoi(context.QueryParam("CartGameId"))
	if err != nil {
		log.Fatal(err)
	}

	canRemoveCartGame := true
	if authOk {
		canRemoveCartGame = db.UserCartGameId(userId, cartGameId, handler.Database) != 0
	}

	if canRemoveCartGame {
		db.RemoveCartGame(cartGameId, handler.Database)
	}

	return context.Redirect(http.StatusMovedPermanently, "/cart")
}

func (handler CartHandler) ClearCartGame(context echo.Context) error {
	time.Sleep(3 * time.Second)

	userId, authOk := Utils{}.GetCookieValueAsInt("Authentication", context)
	cartId, cartIdOk := Utils{}.GetCookieValueAsInt("CartIdentifier", context)

	if authOk {
		db.ClearCartByUserId(userId, handler.Database)
	} else if cartIdOk {
		db.ClearCartByCartId(cartId, handler.Database)
	}

	context.SetCookie(Utils{}.RemoveCookie("CartIdentifier"))

	return context.Redirect(http.StatusMovedPermanently, "/cart")
}
