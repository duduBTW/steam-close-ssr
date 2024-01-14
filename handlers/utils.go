package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Utils struct{}

func (Utils) CreteCookie(key, value string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.Path = "/"

	return cookie
}

func (Utils) RemoveCookie(key string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = ""
	cookie.MaxAge = -1
	cookie.HttpOnly = true
	cookie.Path = "/"

	return cookie
}

func (Utils) GetCookieValueAsInt(key string, context echo.Context) (int, bool) {
	Authentication, err := context.Cookie(key)
	if err != nil {
		return 0, false
	}

	userId, err := strconv.Atoi(Authentication.Value)
	if err != nil {
		return 0, false
	}

	return userId, true
}
