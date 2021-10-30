package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var port string
	var hlsFolder string
	flag.StringVar(&port, "p", "17232", "set hls server port.")
	flag.StringVar(&hlsFolder, "f", "", "set hls folder.")
	flag.Parse()

	if len(hlsFolder) == 0 {
		fmt.Println("Please input hls folder, please type -h to show usage.")
		os.Exit(-1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods: []string{
			http.MethodHead,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		MaxAge: 86400,
	}))

	e.Static("/", hlsFolder)

	go func() {
		if err := e.Start(":" + port); err != nil {
			e.Logger.Errorf(err.Error())
			stop()
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
