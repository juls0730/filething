//go:generate npm --prefix ./ui run generate
package main

import (
	"filething/ui"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	api := e.Group("/api")
	{
		api.GET("/hello", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!!!"})
		})
	}

	e.GET("/*", echo.StaticDirectoryHandler(ui.DistDirFS, false))

	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Logger.Fatal(e.Start(":1323"))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d.html", code)
	file, err := ui.DistDirFS.Open(errorPage)
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
}
