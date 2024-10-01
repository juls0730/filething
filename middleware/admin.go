package middleware

import (
	"filething/models"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func AdminMiddleware(db *bun.DB) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(*models.User)

		if !user.Admin {
			return echo.NewHTTPError(http.StatusForbidden, "You are not an administrator")
		}

		return c.Next()
	}
}
