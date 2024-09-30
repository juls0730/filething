package middleware

import (
	"context"
	"database/sql"
	"filething/models"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

const UserContextKey = "user"

func SessionMiddleware(db *bun.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract sessionToken from the cookie
			cookie, err := c.Cookie("sessionToken")
			if err != nil {
				if err == http.ErrNoCookie {
					return echo.NewHTTPError(http.StatusUnauthorized, "Session token missing")
				}
				return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
			}

			sessionId, err := uuid.Parse(cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
			}

			session := &models.Session{
				ID: sessionId,
			}
			err = db.NewSelect().Model(session).WherePK().Scan(context.Background())

			if err != nil {
				fmt.Println(err)
				if err == sql.ErrNoRows {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid session token")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
			}

			user := &models.User{
				ID: session.UserID,
			}
			err = db.NewSelect().Model(user).Relation("Plan").WherePK().Scan(context.Background())

			if err != nil {
				if err == sql.ErrNoRows {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid session token")
				}
				fmt.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
			}

			// Store the user in the context
			c.Set(UserContextKey, user)

			// Continue to the next handler
			return next(c)
		}
	}
}
