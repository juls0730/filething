// Package ui handles the frontend embedding
package ui

import (
	"embed"

	"github.com/labstack/echo/v4"
)

//go:embed all:.output
var DistDir embed.FS

// DistDirFS contains the embedded dist directory files (without the "dist" prefix)
var DistDirFS = echo.MustSubFS(DistDir, ".output/")
