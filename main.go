package main

import (
	"context"
	"filething/db"
	"filething/middleware"
	"filething/models"
	"filething/routes"
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/uptrace/bun"
)

var initUi func(app *fiber.App)

func main() {
	db.DBConnect()

	db := db.GetDB()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	err = seedPlans(db)
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		JSONEncoder:                  sonic.Marshal,
		JSONDecoder:                  sonic.Unmarshal,
		DisablePreParseMultipartForm: true,
		BodyLimit:                    100 * 1024 * 1024 * 1024,
		StreamRequestBody:            true,
	})

	app.Use(pprof.New())
	// app.Use(compress.New())
	app.Use(cors.New())
	// TODO: make this not a constant pain in my ass
	// app.Use(csrf.New(csrf.Config{
	// 	KeyLookup:      "cookie:_csrf",
	// 	CookieName:     "_csrf",
	// 	CookieSameSite: "Strict",
	// 	Expiration:     time.Hour * 24,
	// 	CookieSecure:   true,
	// 	CookieHTTPOnly: true,
	// }))
	app.Use(helmet.New())

	api := app.Group("/api")
	{
		api.Post("/login", routes.LoginHandler)
		api.Post("/signup", routes.SignupHandler)

		// everything past this needs auth
		api.Use(middleware.SessionMiddleware(db))
		api.Post("/logout", routes.LogoutHandler)
		api.Get("/user", routes.GetUser)

		api.Post("/files/upload*", routes.UploadFile)
		// api.Post("/files/upload/chunked*", routes.UploadFileInChunks)
		api.Get("/files/get/*", routes.GetFiles)
		api.Get("/files/download*", routes.GetFile)
		api.Post("/files/delete*", routes.DeleteFiles)

		admin := api.Group("/admin")
		{
			admin.Use(middleware.AdminMiddleware(db))
			admin.Get("/status", routes.SystemStatus)
			admin.Get("/plans", routes.GetPlans)
			admin.Get("/users", routes.GetUsers)
			admin.Get("/users/:id", routes.GetUser)
			admin.Post("/users/edit/:id", routes.EditUser)
			admin.Post("/users/new", routes.CreateUser)
		}
	}

	// redirects to the proper pages if you are trying to access one that expects you have/dont have an api key
	// this isnt explicitly required, but it provides a better experience than doing this same thing clientside
	app.Use(middleware.AuthCheckMiddleware(db))

	// calls out to a function set by either server.go server_dev.go based on the presence of the dev tag, and hosts
	// either the static files that get embedded into the binary in ui/embed.go or proxies the dev server that gets
	// run in the provided function
	initUi(app)

	routes.AppStartTime = time.Now().UTC()

	if err = app.Listen(":1323"); err != nil && err != http.ErrServerClosed {
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
