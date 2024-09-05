package routes

import (
	"context"
	"filething/models"
	"net/http"
	"regexp"

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

	c.SetCookie(&http.Cookie{
		Name:     "sessionToken",
		Value:    session.ID.String(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful!"})

	// sessionID := uuid.New().String()
	// session := &models.Session{ID: sessionID, UserID: user.ID, ExpiresAt: time.Now().Add(time.Hour * 24)}

	// key := "session:" + session.ID
	// err = client.HSet(ctx, key, session).Err()

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "An unknown error occoured!"})
	// 	return
	// }

	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:  "sessionToken",
	// 	Value: sessionID,
	// 	Path:  "/",
	// })

	// c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
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
	}
	_, err = db.NewInsert().Model(user).Exec(context.Background())

	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "A user with that email or username already exists!"})
	}

	session, err := GenerateSessionToken(db, user.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "An unknown error occoured!"})
	}

	c.SetCookie(&http.Cookie{
		Name:     "sessionToken",
		Value:    session.ID.String(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, map[string]string{"message": "Signup successful!"})

	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:  "sessionToken",
	// 	Value: sessionID,
	// 	Path:  "/",
	// })

	// c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

func GenerateSessionToken(db *bun.DB, userId uuid.UUID) (*models.Session, error) {
	session := &models.Session{
		UserID: userId,
	}

	_, err := db.NewInsert().Model(session).Exec(context.Background())

	return session, err
}
