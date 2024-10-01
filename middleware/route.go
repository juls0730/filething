package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

var unauthenticatedPages = []string{
	"/login",
	"/signup",
	"/",
}

var authenticatedPages = []string{
	"/home",
}

func AuthCheckMiddleware(c fiber.Ctx) error {
	path := c.Path()

	// bypass auth checks for static and dev resources
	if strings.HasPrefix(path, "/_nuxt/") || strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") {
		return c.Next()
	}

	cookie := c.Cookies("sessionToken")
	authenticated := cookie != ""

	if Contains(unauthenticatedPages, path) && authenticated {
		return c.Redirect().To("/home")
	}

	if Contains(authenticatedPages, path) && !authenticated {
		return c.Redirect().To("/login")
	}

	if strings.Contains(path, "/home") && !authenticated {
		return c.Redirect().To("/login")
	}

	if strings.Contains(path, "/admin") && !authenticated {
		return c.Redirect().To("/login")
	}

	return c.Next()
}

func Contains(s []string, element string) bool {
	for _, v := range s {
		if v == element {
			return true
		}
	}
	return false
}
