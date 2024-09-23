//go:generate sh -c "NODE_ENV=production bun --cwd=./ui install"
//go:generate sh -c "NODE_ENV=production bun --bun --cwd=./ui run generate"
package main

import (
	"context"
	"database/sql"
	"filething/middleware"
	"filething/models"
	"filething/routes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var initUi func(e *echo.Echo)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")

	if dbHost == "" || dbName == "" || dbUser == "" || os.Getenv("STORAGE_PATH") == "" {
		panic("Missing database environment variabled!")
	}

	// TODO: retry connection or only connect at the first moment that we need the db
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

	// insert the db into the echo context so it is easily accessible in routes and middleware
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
	e.Use(echoMiddleware.Secure())

	api := e.Group("/api")
	{
		api.POST("/login", routes.LoginHandler)
		api.POST("/signup", routes.SignupHandler)

		// everything past this needs auth
		api.Use(middleware.SessionMiddleware(db))
		api.POST("/logout", routes.LogoutHandler)
		api.GET("/user", routes.GetUser)

		api.POST("/files/upload*", routes.UploadFile)
		api.GET("/files/get/*", routes.GetFiles)
		api.GET("/files/download*", routes.GetFile)
		api.POST("/files/delete*", routes.DeleteFiles)

		admin := api.Group("/admin")
		{
			admin.Use(middleware.AdminMiddleware())
			admin.GET("/system-status", routes.SystemStatus)
			admin.GET("/get-users/:page", routes.GetUsers)
			admin.GET("/get-total-users", routes.GetUsersCount)
		}
	}

	// redirects to the proper pages if you are trying to access one that expects you have/dont have an api key
	// this isnt explicitly required, but it provides a better experience than doing this same thing clientside
	e.Use(middleware.AuthCheckMiddleware)

	// calls out to a function set by either server.go server_dev.go based on the presence of the dev tag, and hosts
	// either the static files that get embedded into the binary in ui/embed.go or proxies the dev server that gets
	// run in the provided function
	initUi(e)

	routes.AppStartTime = time.Now().UTC()

	if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
		fmt.Println("Error starting HTTP server:", err)
	}
}

// creates tables in the db if they dont already exist
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

// seeds the storage plans into the database
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

	// TODO: read this from a config
	plans := []models.Plan{
		{MaxStorage: 10 * 1024 * 1024 * 1024},  // 10GB
		{MaxStorage: 50 * 1024 * 1024 * 1024},  // 50GB
		{MaxStorage: 100 * 1024 * 1024 * 1024}, // 100GB
	}

	_, err = db.NewInsert().Model(&plans).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to seed plans: %w", err)
	}

	return nil
}
