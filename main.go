//go:generate bun --cwd=./ui install
//go:generate bun --bun --cwd=./ui run generate
package main

import (
	"context"
	"database/sql"
	"filething/middleware"
	"filething/models"
	"filething/routes"
	"filething/ui"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")

	if dbHost == "" || dbName == "" || dbUser == "" || os.Getenv("STORAGE_PATH") == "" {
		panic("Missing database environment variabled!")
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?dial_timeout=10s&sslmode=disable", dbUser, dbPasswd, dbHost, dbName)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))
	db := bun.NewDB(sqldb, pgdialect.New())

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	err = seedPlans(db)
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	e.Use(echoMiddleware.Gzip())
	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.CSRFWithConfig(echoMiddleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	}))

	api := e.Group("/api")
	{
		api.POST("/login", routes.LoginHandler)
		api.POST("/signup", routes.SignupHandler)

		// everything past this needs auth
		api.Use(middleware.SessionMiddleware(db))
		api.GET("/user", routes.GetUser)
		api.GET("/user/usage", routes.GetUsage)

		api.POST("/upload*", routes.UploadFile)
		api.GET("/files*", routes.GetFiles)
	}

	// redirects to the proper pages if you are trying to access one that expects you have/dont have an api key
	// this isnt explicitly required, but it provides a better experience than doing this same thing clientside
	e.Use(middleware.AuthCheckMiddleware)

	e.GET("/*", echo.StaticDirectoryHandler(ui.DistDirFS, false))

	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Logger.Fatal(e.Start(":1323"))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusNotFound {
		path := c.Request().URL.Path

		if !strings.HasPrefix(path, "/api") {
			file, err := ui.DistDirFS.Open("404.html")
			if err != nil {
				c.Logger().Error(err)
			}

			fileInfo, err := file.Stat()
			if err != nil {
				c.Logger().Error(err)
			}

			fileBuf := make([]byte, fileInfo.Size())
			_, err = file.Read(fileBuf)
			defer func() {
				if err := file.Close(); err != nil {
					panic(err)
				}
			}()
			if err != nil {
				c.Logger().Error(err)
				panic(err)
			}

			c.HTML(http.StatusNotFound, string(fileBuf))
			return
		}
	}

	c.Echo().DefaultHTTPErrorHandler(err, c)
}

func createSchema(db *bun.DB) error {
	models := []interface{}{
		(*models.User)(nil),
		(*models.Session)(nil),
		(*models.Plan)(nil),
	}

	ctx := context.Background()
	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedPlans(db *bun.DB) error {
	ctx := context.Background()
	count, err := db.NewSelect().Model((*models.Plan)(nil)).Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count plans: %w", err)
	}

	// If the table is not empty, no need to seed
	if count > 0 {
		return nil
	}

	plans := []models.Plan{
		{MaxStorage: 10 * 1024 * 1024 * 1024},  // 10GB
		{MaxStorage: 50 * 1024 * 1024 * 1024},  // 50GB
		{MaxStorage: 100 * 1024 * 1024 * 1024}, // 100GB
	}

	_, err = db.NewInsert().Model(&plans).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to seed plans: %w", err)
	}

	log.Println("Successfully seeded the plans table")
	return nil
}
