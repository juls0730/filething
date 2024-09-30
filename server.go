//go:build !dev && !ssr
// +build !dev,!ssr

package main

import (
	"filething/ui"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type embeddedFS struct {
	baseFS fs.FS
	prefix string
}

func (fs *embeddedFS) Open(name string) (fs.File, error) {
	// Prepend the prefix to the requested file name
	publicPath := filepath.Join(fs.prefix, name)
	return fs.baseFS.Open(publicPath)
}

var publicFS = &embeddedFS{
	baseFS: ui.DistDirFS,
	prefix: "public/",
}

func init() {
	initUi = func(e *echo.Echo) {
		// e.GET("/*", echo.StaticDirectoryHandler(publicFS, false))

		e.HTTPErrorHandler = customHTTPErrorHandler
	}
}

// Custom Error handling since Nuxt relies on the 404 page for dynamic pages we still want api routes to use the default
// error handling built into echo
func customHTTPErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusNotFound {
		path := c.Request().URL.Path

		if !strings.HasPrefix(path, "/api") {
			file, err := publicFS.Open("404.html")
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
