package routes

import (
	"context"
	"filething/models"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c echo.Context) error {
	var loginData models.LoginData

	if err := c.Bind(&loginData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An unknown error occoured!"})
	}

	if loginData.UsernameOrEmail == "" || loginData.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "A password, and username or email are required!"})
	}

	db := c.Get("db").(*bun.DB)

	user := new(models.User)
	err := db.NewSelect().Model(user).Where("email = ?", loginData.UsernameOrEmail).Relation("Plan").Scan(context.Background())
	if err != nil {
		err := db.NewSelect().Model(user).Where("username = ?", loginData.UsernameOrEmail).Relation("Plan").Scan(context.Background())
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "User with that username or email not found!"})
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An unknown error occoured!"})
	}

	expiration := time.Now().Add(time.Hour * 24 * 365 * 100)

	c.SetCookie(&http.Cookie{
		Name:     "sessionToken",
		Value:    session.ID.String(),
		SameSite: http.SameSiteStrictMode,
		Expires:  expiration,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, user)
}

var firstUserCreated *bool

func SignupHandler(c echo.Context) error {
	var signupData models.SignupData

	if err := c.Bind(&signupData); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An unknown error occoured!"})
	}

	if signupData.Username == "" || signupData.Password == "" || signupData.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "A password, username and email are required!"})
	}

	// if email is not valid
	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(signupData.Email) {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "A valid email is required!"})
	}

	db := c.Get("db").(*bun.DB)

	hash, err := bcrypt.GenerateFromPassword([]byte(signupData.Password), 12)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An unknown error occoured!"})
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
		return c.JSON(http.StatusConflict, map[string]string{"message": "A user with that email or username already exists!"})
	}

	if !*firstUserCreated {
		*firstUserCreated = true
	}

	err = db.NewSelect().Model(user).WherePK().Relation("Plan").Scan(context.Background())
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "An unknown error occoured!"})
	}

	err = os.Mkdir(fmt.Sprintf("%s/%s", os.Getenv("STORAGE_PATH"), user.ID), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	session, err := GenerateSessionToken(db, user.ID)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An unknown error occoured!"})
	}

	expiration := time.Now().Add(time.Hour * 24 * 365 * 100)

	c.SetCookie(&http.Cookie{
		Name:     "sessionToken",
		Value:    session.ID.String(),
		SameSite: http.SameSiteStrictMode,
		Expires:  expiration,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, user)
}

func GenerateSessionToken(db *bun.DB, userId uuid.UUID) (*models.Session, error) {
	session := &models.Session{
		UserID: userId,
	}

	_, err := db.NewInsert().Model(session).Exec(context.Background())

	return session, err
}

func GetUser(c echo.Context) error {
	if c.Param("id") == "" {
		user := c.Get("user")
		if user == nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
		}

		basePath := fmt.Sprintf("%s/%s/", os.Getenv("STORAGE_PATH"), user.(*models.User).ID)
		storageUsage, err := calculateStorageUsage(basePath)
		if err != nil {
			return err
		}

		user.(*models.User).Usage = storageUsage

		return c.JSON(http.StatusOK, user.(*models.User))
	} else {
		// get a user from the db using the id parameter, this *should* only be used for admin since /api/admin/users/:id has
		// a middleware that checks if the user is an admin, and it should be impossible to pass a param to this endpoint if it isnt that route
		db := c.Get("db").(*bun.DB)

		user := new(models.User)
		userId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "An unknown error occoured!"})
		}

		user.ID = userId

		err = db.NewSelect().Model(user).WherePK().Relation("Plan").Scan(context.Background())
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
		}

		basePath := fmt.Sprintf("%s/%s/", os.Getenv("STORAGE_PATH"), user.ID)
		storageUsage, err := calculateStorageUsage(basePath)
		if err != nil {
			return err
		}

		user.Usage = storageUsage

		return c.JSON(http.StatusOK, user)
	}
}

func LogoutHandler(c echo.Context) error {
	db := c.Get("db").(*bun.DB)

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
	_, err = db.NewDelete().Model(session).WherePK().Exec(context.Background())

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An unknown error occoured!"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Succesfully logged out"})
}
