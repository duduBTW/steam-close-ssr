package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"database/sql"

	"github.com/duduBTW/steam-ssr/db"
	"github.com/duduBTW/steam-ssr/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Running")

	db, err := sql.Open("sqlite3", "steam.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	echoInstance := echo.New()
	echoInstance.Use(authMiddleware(db))
	echoInstance.Use(cartMiddleware(db))

	echoInstance.Static("/front", "./front/dist")
	echoInstance.Static("/static", "./static")

	// Login
	echoInstance.GET("/login", handlers.Login{}.Page)
	echoInstance.GET("/register", handlers.Register{}.Page)
	echoInstance.POST("/api/register", handlers.Register{Database: db}.Register)
	echoInstance.POST("/api/login", handlers.Login{Database: db}.Login)
	echoInstance.GET("/api/logout", handlers.Login{Database: db}.Logout)

	// Home
	echoInstance.GET("/", handlers.HomeHandler{Database: db}.HomePage)

	// Cart
	echoInstance.GET("/cart", handlers.CartHandler{Database: db}.Page)
	echoInstance.GET("/api/cart", handlers.CartHandler{Database: db}.InsertCart)
	echoInstance.GET("/api/cart/delete", handlers.CartHandler{Database: db}.RemoveCartGame)
	echoInstance.GET("/api/cart/clear", handlers.CartHandler{Database: db}.ClearCartGame)

	// Game
	echoInstance.GET("/game/:gameId", handlers.GameHandler{Database: db}.GamePage)
	echoInstance.GET("/api/game/wishlist", handlers.WishlistHandler{Database: db}.WishlistButton)
	echoInstance.GET("/api/game/search/fast", handlers.GameHandler{Database: db}.FastSearch)

	// Logger
	echoInstance.Logger.Fatal(echoInstance.Start("127.0.0.1:3001"))
	echoInstance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
}

func cartMiddleware(Database *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId, successAuth := handlers.Utils{}.GetCookieValueAsInt("Authentication", c)
			cartId, successCart := handlers.Utils{}.GetCookieValueAsInt("CartIdentifier", c)
			if !successAuth && !successCart {
				return next(c)
			}

			setContextValue("CartCount", db.CountCartEntries(cartId, userId, Database), c)

			return next(c)
		}
	}
}

func authMiddleware(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			setContextValue("CurrentUrl", c.Request().URL.Path, c)
			AuthenticationCookie, err := c.Cookie("Authentication")
			if err != nil {
				return nextNotAuthenticated(next, c)
			}

			userId, err := strconv.Atoi(AuthenticationCookie.Value)
			if err != nil {
				return nextNotAuthenticated(next, c)
			}

			user := handlers.FindUserById(userId, db)
			setContextValue("Authenticated", true, c)
			setContextValue("DisplayName", user.DisplayName, c)
			setContextValue("ProfilePicture", user.ProfilePicture, c)
			setContextValue("UserId", user.UserId, c)

			return next(c)
		}
	}
}

func nextNotAuthenticated(next echo.HandlerFunc, c echo.Context) error {
	setContextValue("Authenticated", false, c)
	return next(c)
}

func setContextValue(key string, value any, c echo.Context) {
	c.Set(key, value)
	c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), key, value)))
}
