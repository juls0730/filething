package middleware

import (
	"context"
	"database/sql"
	"filething/models"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const UserContextKey = "user"

func SessionMiddleware(db *bun.DB) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		// Extract session token from the cookie
		sessionToken := c.Cookies("sessionToken")
		if sessionToken == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Session token missing"})
		}

		// Parse session ID
		sessionId, err := uuid.Parse(sessionToken)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid session token"})
		}

		// Fetch session from database
		session := &models.Session{
			ID: sessionId,
		}
		err = db.NewSelect().Model(session).WherePK().Scan(context.Background())

		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid session token"})
			}
			fmt.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Database error"})
		}

		user := &models.User{
			ID: session.UserID,
		}
		err = db.NewSelect().Model(user).Relation("Plan").WherePK().Scan(context.Background())

		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid session token"})
			}
			fmt.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Database error"})
		}

		c.Locals("user", user)

		return c.Next()
	}
}
