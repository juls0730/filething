package middleware

import (
	"context"
	"filething/models"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var unauthenticatedPages = []string{
	"/login",
	"/signup",
	"/",
}

var authenticatedPages = []string{
	"/home",
	"/admin",
}

func AuthCheckMiddleware(db *bun.DB) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		path := c.Path()

		// bypass auth checks for static and dev resources
		if strings.HasPrefix(path, "/_nuxt/") || strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") {
			return c.Next()
		}

		var authenticated bool = true
		sessionToken := c.Cookies("sessionToken")
		if sessionToken == "" {
			authenticated = false
		}

		// Parse session ID
		sessionId, err := uuid.Parse(sessionToken)
		if err != nil {
			authenticated = false
		}

		// Fetch session from database
		session := &models.Session{
			ID: sessionId,
		}
		err = db.NewSelect().Model(session).WherePK().Scan(context.Background())

		if err != nil {
			authenticated = false
		}

		if Contains(unauthenticatedPages, path) && authenticated {
			fmt.Println("unauthenticated page", path, authenticated)
			return c.Redirect().To("/home")
		}

		if Contains(authenticatedPages, path) && !authenticated {
			fmt.Println("authenticated page", path, authenticated)
			return c.Redirect().To("/login")
		}

		if strings.Contains(path, "/home") && !authenticated {
			fmt.Println("home page", path, authenticated)
			return c.Redirect().To("/login")
		}

		if strings.Contains(path, "/admin") && !authenticated {
			fmt.Println("admin page", path, authenticated)
			return c.Redirect().To("/login")
		}

		return c.Next()
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
