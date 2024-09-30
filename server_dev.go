//go:build dev
// +build dev

package main

import (
	"net/url"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func init() {
	initUi = func(e *echo.Echo) {
		spawnProcess("bun", []string{"--cwd=ui", "run", "dev"}, e)

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
