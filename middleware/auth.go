package middleware

import (
	"context"
	"database/sql"
	"filething/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

// import (
// 	"database/sql"
// 	"net/http"

// 	"github.com/go-pg/pg/v10"

// 	"github.com/labstack/echo/v4"
// )

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

			sessionToken := cookie.Value

			// Query the session and user data from PostgreSQL
			session := new(models.Session)
			err = db.NewSelect().Model(session).Relation("User").WherePK(sessionToken).Scan(context.Background())

			if err != nil {
				if err == sql.ErrNoRows {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid session token")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
			}

			user := session.User

			// Store the user in the context
			c.Set(UserContextKey, user)

			// Continue to the next handler
			return next(c)
		}
	}
}
