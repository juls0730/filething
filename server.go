//go:build !dev
// +build !dev

package main

import (
	"filething/ui"

	"github.com/labstack/echo/v4"
)

func init() {
	initUi = func(e *echo.Echo) {
		e.GET("/*", echo.StaticDirectoryHandler(ui.DistDirFS, false))
	}
}
