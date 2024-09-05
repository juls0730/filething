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
			sessionId, err := uuid.Parse(sessionToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
			}

			session := &models.Session{
				ID: sessionId,
			}
			err = db.NewSelect().Model(session).Relation("User").WherePK().Scan(context.Background())

			if err != nil {
				fmt.Println(err)
				if err == sql.ErrNoRows {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid session token")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
			}

			user := &session.User

			// Store the user in the context
			c.Set(UserContextKey, user)

			// Continue to the next handler
			return next(c)
		}
	}
}
