//go:build !dev && !ssr
// +build !dev,!ssr

package main

import (
	"filething/ui"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type embeddedFS struct {
	baseFS fs.FS
	prefix string
}

func (fs *embeddedFS) Open(name string) (fs.File, error) {
	// Prepend the prefix to the requested file name
	publicPath := filepath.Join(fs.prefix, name)
	fmt.Println("Reading file:", publicPath)
	file, err := fs.baseFS.Open(publicPath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %s", publicPath)
	}

	fmt.Println("File found:", publicPath, file)

	return file, err
}

var publicFS = &embeddedFS{
	baseFS: ui.DistDir,
	prefix: "public/",
}

func init() {
	initUi = func(app *fiber.App) {
		app.Get("/*", static.New("", static.Config{
			FS: publicFS,
		}))

		app.Use(func(c fiber.Ctx) error {
			err := c.Next()
			if err == nil {
				return nil
			}

			if fiber.ErrNotFound == err {
				path := c.Path()
				if !strings.HasPrefix(path, "/api") {
					file, err := publicFS.Open("404.html")
					if err != nil {
						c.App().Server().Logger.Printf("Error opening 404.html: %s", err)
						return err
					}
					defer file.Close()

					fileInfo, err := file.Stat()
					if err != nil {
						c.App().Server().Logger.Printf("An error occurred while getting the file info: %s", err)
						return err
					}

					fileBuf := make([]byte, fileInfo.Size())
					_, err = file.Read(fileBuf)
					if err != nil {
						c.App().Server().Logger.Printf("An error occurred while reading the file: %s", err)
						return err
					}

					return c.Status(fiber.StatusNotFound).SendString(string(fileBuf))
				}
			}
			return err
		})
	}
}
