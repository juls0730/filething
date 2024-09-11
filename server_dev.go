//go:build dev
// +build dev

package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func init() {
	initUi = func(e *echo.Echo) {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			cmd := exec.Command("bun", "--cwd=ui", "run", "dev")
			cmd.Stderr = os.Stderr

			// use a preocess group since otherwise the node processes spawned by bun wont die
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

			err := cmd.Start()
			if err != nil {
				if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
					fmt.Println("Error sending SIGTERM to Nuxt dev server group:", err)
				}

				fmt.Println("Error starting Nuxt dev server:", err)
				return
			}

			go func() {
				<-shutdown

				if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
					fmt.Println("Error sending SIGTERM to Nuxt dev server group:", err)
				}
			}()

			if err := cmd.Wait(); err != nil {
				fmt.Println("Error waiting for Nuxt dev server to exit:", err)
			}

			fmt.Println("Nuxt dev server stopped")

			if err := e.Shutdown(context.Background()); err != nil {
				fmt.Println("Error shutting down HTTP server:", err)
			}
		}()

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
