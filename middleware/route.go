package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var unauthenticatedPages = []string{
	"/login",
	"/signup",
	"/",
}

var authenticatedPages = []string{
	"/home",
}

func AuthCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path

		// bypass auth checks for static and dev resources
		if strings.HasPrefix(path, "/_nuxt/") || strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") {
			return next(c)
		}

		_, cookieErr := c.Cookie("sessionToken")
		authenticated := cookieErr == nil

		if Contains(unauthenticatedPages, path) && authenticated {
			return c.Redirect(http.StatusFound, "/home")
		}

		if Contains(authenticatedPages, path) && !authenticated {
			return c.Redirect(http.StatusFound, "/login")
		}

		if strings.Contains(path, "/home") && !authenticated {
			return c.Redirect(http.StatusFound, "/login")
		}

		if strings.Contains(path, "/admin") && !authenticated {
			return c.Redirect(http.StatusFound, "/login")
		}

		return next(c)
	}
}

func Contains(s []string, element string) bool {
	for _, v := range s {
		if v == element {
			return true
		}
	}
	return false
}
