//go:generate bun --cwd=./ui install
//go:generate bun --cwd=./ui run generate
package main

import (
	"context"
	"database/sql"
	"filething/models"
	"filething/routes"
	"filething/ui"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")

	if dbHost == "" || dbName == "" || dbUser == "" {
		panic("Missing database environment variabled!")
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?dial_timeout=10s&sslmode=disable", dbUser, dbPasswd, dbHost, dbName)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))
	db := bun.NewDB(sqldb, pgdialect.New())

	err := createSchema(db)
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

	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
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
		api.GET("/hello", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!!!"})
		})
	}

	e.GET("/*", echo.StaticDirectoryHandler(ui.DistDirFS, false))

	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Logger.Fatal(e.Start(":1323"))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)

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
