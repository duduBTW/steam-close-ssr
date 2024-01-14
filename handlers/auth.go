package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/duduBTW/steam-ssr/data"
	"github.com/duduBTW/steam-ssr/db"
	"github.com/duduBTW/steam-ssr/presentation/auth"
	"github.com/labstack/echo/v4"
)

type Login struct {
	Database *sql.DB
}

func (h Login) Page(c echo.Context) error {
	return auth.LoginPage().Render(c.Request().Context(), c.Response())
}

func (handler Login) Login(c echo.Context) error {
	user := FindUserByUserName(c.FormValue("UserName"), handler.Database)

	authCookie := Utils{}.CreteCookie("Authentication", strconv.Itoa(user.UserId))
	c.SetCookie(authCookie)

	cartId, cartIdOk := Utils{}.GetCookieValueAsInt("CartIdentifier", c)
	if cartIdOk {
		db.VinculateUserToCart(user.UserId, cartId, handler.Database)
		c.SetCookie(Utils{}.RemoveCookie("CartIdentifier"))
	}

	redirectUrl := c.QueryParam("redirectUrl")
	url := "/"
	if redirectUrl != "" {
		url = redirectUrl
	}

	return c.Redirect(http.StatusMovedPermanently, url)
}

func (handler Login) Logout(c echo.Context) error {
	authCookie := Utils{}.RemoveCookie("Authentication")
	c.SetCookie(authCookie)

	redirectUrl := c.QueryParam("redirectUrl")
	url := "/"
	if redirectUrl != "" {
		url = redirectUrl
	}

	return c.Redirect(http.StatusMovedPermanently, url)
}

type Register struct {
	Database *sql.DB
}

func (h Register) Page(c echo.Context) error {
	return auth.RegisterPage().Render(c.Request().Context(), c.Response())
}

func (handler Register) Register(c echo.Context) error {
	InsetUser(data.User{
		UserName:       c.FormValue("UserName"),
		ProfilePicture: c.FormValue("ProfilePicture"),
		DisplayName:    c.FormValue("DisplayName"),
	}, handler.Database)

	return c.Redirect(http.StatusMovedPermanently, "/login")
}

func InsetUser(user data.User, db *sql.DB) int {
	query := `
		INSERT INTO User
			(UserName, ProfilePicture, DisplayName) 
		VALUES 
			($1, $2, $3)
		RETURNING
			UserId
	`

	var UserId int
	err := db.QueryRow(query, user.UserName, user.ProfilePicture, user.DisplayName).Scan(&UserId)
	if err != nil {
		log.Fatal(err)
	}

	return UserId
}

func FindUserByUserName(userName string, db *sql.DB) data.User {
	query := `
		SELECT 
			UserId,
			UserName,
			ProfilePicture,
			DisplayName
		FROM
			User
		WHERE	
			UserName = $1
	`

	var user data.User
	err := db.QueryRow(query, userName).Scan(&user.UserId, &user.UserName, &user.ProfilePicture, &user.DisplayName)
	if err != nil {
		log.Fatal(err)
	}

	return user
}

func FindUserById(userId int, db *sql.DB) data.User {
	query := `
		SELECT 
			UserId,
			UserName,
			ProfilePicture,
			DisplayName
		FROM
			User
		WHERE	
		UserId = $1
	`

	var user data.User
	err := db.QueryRow(query, userId).Scan(&user.UserId, &user.UserName, &user.ProfilePicture, &user.DisplayName)
	if err != nil {
		log.Fatal(err)
	}

	return user
}
