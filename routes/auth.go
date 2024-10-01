package routes

import (
	"context"
	"filething/db"
	"filething/models"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c fiber.Ctx) error {
	loginData := new(models.LoginData)

	if err := c.Bind().JSON(loginData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "An unknown error occoured!"})
	}

	// Validate required fields
	if loginData.UsernameOrEmail == "" || loginData.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "A Username/Email and Password are required"})
	}

	db := db.GetDB()

	user := new(models.User)
	var err error

	// Try both username and email for login attempts
	if err = db.NewSelect().Model(user).Where("email = ?", loginData.UsernameOrEmail).Relation("Plan").Scan(context.Background()); err != nil {
		err = db.NewSelect().Model(user).Where("username = ?", loginData.UsernameOrEmail).Relation("Plan").Scan(context.Background())
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
	}

	basePath := fmt.Sprintf("%s/%s/", os.Getenv("STORAGE_PATH"), user.ID)
	storageUsage, err := calculateStorageUsage(basePath)
	if err != nil {
		return err
	}

	user.Usage = storageUsage

	session, err := GenerateSessionToken(db, user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	expiration := time.Now().Add(time.Hour * 24 * 365 * 100)

	c.Cookie(&fiber.Cookie{
		Name:     "sessionToken",
		Value:    session.ID.String(),
		SameSite: "Strict",
		Expires:  expiration,
		Path:     "/",
	})

	// Return user data with status code 200 (OK)
	return c.JSON(user)
}

var firstUserCreated *bool

func SignupHandler(c fiber.Ctx) error {
	signupData := new(models.SignupData)

	if err := c.Bind().JSON(signupData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "An unknown error occoured!"})
	}

	if signupData.Username == "" || signupData.Password == "" || signupData.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "A password, username and email are required!"})
	}

	// if email is not valid
	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(signupData.Email) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "A valid email is required!"})
	}

	db := db.GetDB()

	hash, err := bcrypt.GenerateFromPassword([]byte(signupData.Password), 12)
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "An unknown error occoured!"})
	}

	if firstUserCreated == nil {
		count, err := db.NewSelect().Model((*models.User)(nil)).Count(context.Background())
		if err != nil {
			return fmt.Errorf("failed to count plans: %w", err)
		}

		firstUserCreated = new(bool)
		*firstUserCreated = count != 0
	}

	user := &models.User{
		Username:     signupData.Username,
		Email:        signupData.Email,
		PasswordHash: string(hash),
		PlanID:       1, // basic 10GB plan
		Admin:        !*firstUserCreated,
	}
	_, err = db.NewInsert().Model(user).Exec(context.Background())

	if err != nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"message": "A user with that email or username already exists!"})
	}

	if !*firstUserCreated {
		*firstUserCreated = true
	}

	err = db.NewSelect().Model(user).WherePK().Relation("Plan").Scan(context.Background())
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "An unknown error occoured!"})
	}

	err = os.Mkdir(fmt.Sprintf("%s/%s", os.Getenv("STORAGE_PATH"), user.ID), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	session, err := GenerateSessionToken(db, user.ID)

	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "An unknown error occoured!"})
	}

	expiration := time.Now().Add(time.Hour * 24 * 365 * 100)

	c.Cookie(&fiber.Cookie{
		Name:     "sessionToken",
		Value:    session.ID.String(),
		SameSite: "Strict",
		Expires:  expiration,
		Path:     "/",
	})

	return c.Status(http.StatusOK).JSON(user)
}

func GenerateSessionToken(db *bun.DB, userId uuid.UUID) (*models.Session, error) {
	session := &models.Session{
		UserID: userId,
	}

	_, err := db.NewInsert().Model(session).Exec(context.Background())

	return session, err
}

func GetUser(c fiber.Ctx) error {
	if c.Params("id") == "" {
		user := c.Locals("user")
		if user == nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}

		basePath := fmt.Sprintf("%s/%s/", os.Getenv("STORAGE_PATH"), user.(*models.User).ID)
		storageUsage, err := calculateStorageUsage(basePath)
		if err != nil {
			return err
		}

		user.(*models.User).Usage = storageUsage

		return c.Status(http.StatusOK).JSON(user.(*models.User))
	} else {
		// get a user from the db using the id parameter, this *should* only be used for admin since /api/admin/users/:id has
		// a middleware that checks if the user is an admin, and it should be impossible to pass a param to this endpoint if it isnt that route
		db := db.GetDB()

		user := new(models.User)
		userId, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "An unknown error occoured!"})
		}

		user.ID = userId

		err = db.NewSelect().Model(user).WherePK().Relation("Plan").Scan(context.Background())
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}

		basePath := fmt.Sprintf("%s/%s/", os.Getenv("STORAGE_PATH"), user.ID)
		storageUsage, err := calculateStorageUsage(basePath)
		if err != nil {
			return err
		}

		user.Usage = storageUsage

		return c.Status(http.StatusOK).JSON(user)
	}
}

func LogoutHandler(c fiber.Ctx) error {
	db := db.GetDB()

	cookie := c.Cookies("sessionToken")
	if cookie == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Session token missing"})
	}

	sessionId, err := uuid.Parse(cookie)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Bad request"})
	}

	session := &models.Session{
		ID: sessionId,
	}
	_, err = db.NewDelete().Model(session).WherePK().Exec(context.Background())

	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "An unknown error occoured!"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Succesfully logged out"})
}
