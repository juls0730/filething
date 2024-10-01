//go:build dev
// +build dev

package main

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
	fastProxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

func init() {
	initUi = func(app *fiber.App) {
		if !fiber.IsChild() {
			spawnProcess("bun", []string{"--cwd=ui", "run", "dev"}, app)
		}

		target := "localhost:3000"
		app.All("/*", func(c fiber.Ctx) error {
			path := c.Path()
			if strings.HasPrefix(path, "/api") {
				return c.Next()
			}

			request := c.Request().URI()
			if string(request.RequestURI()) == "/_nuxt/" {
				return proxyWebSocket(c, target)
			}

			return proxy.Do(c, "http://"+target+string(request.RequestURI()))
		})
	}
}

var proxyServer *fastProxy.ReverseProxy

func proxyWebSocket(c fiber.Ctx, target string) error {
	path := c.Path()
	// proxyServer, err := fastProxy.NewWSReverseProxyWith(
	// 	fastProxy.WithURL_OptionWS("ws://localhost:3000"+path),
	// 	fastProxy.WithDynamicPath_OptionWS(true, fastProxy.DefaultOverrideHeader),
	// )
	if proxyServer == nil {
		proxyServer, err = fastProxy.NewWSReverseProxyWith(
			fastProxy.WithURL_OptionWS("ws://localhost:3000"+path),
			fastProxy.WithDynamicPath_OptionWS(true, fastProxy.DefaultOverrideHeader),
		)
		if err != nil {
			panic(err)
		}
	}
	proxyServer.ServeHTTP(c.Context())
	return nil
}
