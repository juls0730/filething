package middleware

import (
	"filething/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*models.User)

			if !user.Admin {
				return echo.NewHTTPError(http.StatusForbidden, "You are not an administrator")
			}

			return next(c)
		}
	}
}
