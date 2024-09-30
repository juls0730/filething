//go:build ssr
// +build ssr

package main

import (
	"embed"
	"filething/ui"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func init() {
	initUi = func(e *echo.Echo) {
		tmpDir, err := os.MkdirTemp("", "filething-ssr")
		if err != nil {
			panic(err)
		}

		err = copyEmbeddedFiles(ui.DistDir, ".output", tmpDir)
		if err != nil {
			panic(err)
		}

		path := filepath.Join(tmpDir, "server/index.mjs")
		spawnProcess("node", []string{path}, e)

		target := "localhost:3000"
		e.Group("/*").Use(echoMiddleware.ProxyWithConfig(echoMiddleware.ProxyConfig{
			Balancer: echoMiddleware.NewRoundRobinBalancer([]*echoMiddleware.ProxyTarget{
				{URL: &url.URL{
					Scheme: "http",
					Host:   target,
				}},
			}),
		}))
	}
}

func copyEmbeddedFiles(fs embed.FS, sourcePath string, targetDir string) error {
	entries, err := fs.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourceFile := filepath.Join(sourcePath, entry.Name())
		destFile := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			os.Mkdir(destFile, 0755)
			err := copyEmbeddedFiles(fs, sourceFile, destFile)
			if err != nil {
				return err
			}
		} else {
			source, err := fs.Open(sourceFile)
			if err != nil {
				return err
			}
			defer source.Close()

			dest, err := os.Create(destFile)
			if err != nil {
				return err
			}
			defer dest.Close()

			_, err = io.Copy(dest, source)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
