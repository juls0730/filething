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
	err := db.NewSelect().Model(user).Where("email = ?", loginData.UsernameOrEmail).Scan(context.Background())
	if err != nil {
		err := db.NewSelect().Model(user).Where("username = ?", loginData.UsernameOrEmail).Scan(context.Background())
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "User with that username or email not found!"})
		}
	}

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

	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful!"})
}

func SignupHandler(c echo.Context) error {
	var signupData models.SignupData

	if err := c.Bind(&signupData); err != nil {
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An unknown error occoured!"})
	}

	user := &models.User{
		Username:     signupData.Username,
		Email:        signupData.Email,
		PasswordHash: string(hash),
		PlanID:       1, // basic 10GB plan
	}
	_, err = db.NewInsert().Model(user).Exec(context.Background())

	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "A user with that email or username already exists!"})
	}

	err = os.Mkdir(fmt.Sprintf("%s/%s", os.Getenv("STORAGE_PATH"), user.ID), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An unknown error occoured!"})
	}

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

	return c.JSON(http.StatusOK, map[string]string{"message": "Signup successful!"})
}

func GenerateSessionToken(db *bun.DB, userId uuid.UUID) (*models.Session, error) {
	session := &models.Session{
		UserID: userId,
	}

	_, err := db.NewInsert().Model(session).Exec(context.Background())

	return session, err
}

func GetUser(c echo.Context) error {
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
}
