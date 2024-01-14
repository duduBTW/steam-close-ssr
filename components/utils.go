package components

import (
	"context"
	"strconv"
)

func GetUserName(c context.Context) string {
	name, ok := c.Value("DisplayName").(string)
	if !ok {
		return ""
	}

	return name
}

func getProfilePicture(c context.Context) string {
	profilePicture, ok := c.Value("ProfilePicture").(string)
	if !ok {
		return ""
	}

	return profilePicture
}

func isAuthenticated(c context.Context) bool {
	authenticated, ok := c.Value("Authenticated").(bool)
	if !ok {
		return false
	}

	return authenticated
}

func getCurrentUrl(c context.Context) string {
	currentUrl, ok := c.Value("CurrentUrl").(string)
	if !ok {
		return ""
	}

	return currentUrl
}

func generateRedirectUrl(c context.Context) string {
	currentUrl := getCurrentUrl(c)
	if currentUrl == "" {
		return "/"
	}

	return "?redirectUrl=" + getCurrentUrl(c)
}

func getCartCount(c context.Context) string {
	cartCount, ok := c.Value("CartCount").(int)
	if !ok {
		return ""
	}

	return strconv.Itoa(cartCount)
}

func hasCartCount(c context.Context) bool {
	cartCount, ok := c.Value("CartCount").(int)
	if !ok {
		return false
	}

	return cartCount > 0
}
