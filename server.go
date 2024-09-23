//go:build !dev
// +build !dev

package main

import (
	"filething/ui"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func init() {
	initUi = func(e *echo.Echo) {
		e.GET("/*", echo.StaticDirectoryHandler(ui.DistDirFS, false))

		e.HTTPErrorHandler = customHTTPErrorHandler
	}
}

// Custom Error handling since Nuxt relies on the 404 page for dynamic pages we still want api routes to use the default
// error handling built into echo
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
