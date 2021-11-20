package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	botsFile := flag.String("f", findBotsFile(), "Use the specified .yaml file")
	listenAddr := flag.String("l", "localhost:8443", "Listen address")
	strict := flag.Bool("s", false, "Strict mode - blocks requests not coming from Telegram")

	flag.Parse()

	if *botsFile == "" {
		fmt.Println("Default bots.yml file not found. Please create it or specify one with -f.")
		os.Exit(1)
	}

	botsData, err := readBotsFile(*botsFile)
	if err != nil {
		fmt.Printf("Failed to load bots file: %v", err)
		os.Exit(1)
	}

	// calculate dynamic paddings based on the longest values
	var namePd, endpointPd int
	for _, b := range botsData.Bots {
		if namePd < len(b.Name) {
			namePd = len(b.Name)
		}
		if endpointPd < len(b.Endpoint) {
			endpointPd = len(b.Endpoint)
		}
	}

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "[${time_custom}] ${remote_ip} -> ${path} [${status}] (${latency_human})\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	if *strict {
		fmt.Println("Strict mode enabled!")
		e.Use(telegramIPMiddleware)
	}

	fmt.Println("Loading routes...")

	for _, b := range botsData.Bots {
		b.NormalizeEndpoint()
		fmt.Printf("    %"+strconv.Itoa(namePd)+"s: %-"+strconv.Itoa(endpointPd)+"s -> %s\n", b.Name, b.Endpoint, b.Path)
		registerRoute(e, b.Endpoint, b.Path)
	}

	go func() {
		err := e.Start(*listenAddr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server start failed:", err)
			os.Exit(1)
		}
	}()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)

	<-exitChan
	fmt.Println("Interrupt received, stopping webserver...")
	// TODO: implement pressing Ctrl + C again for forcing exit
	if err = e.Shutdown(context.Background()); err != nil {
		fmt.Println("Server shutdown failed:", err)
	}
}

func registerRoute(e *echo.Echo, endpoint string, path string) {
	e.POST(endpoint, func(c echo.Context) error {
		resp, err := http.Post(path, c.Request().Header.Get("Content-Type"), c.Request().Body)
		if err != nil {
			// TODO: improve logging
			fmt.Println(err)
			return c.NoContent(http.StatusGatewayTimeout)
		}
		return c.NoContent(resp.StatusCode)
	})
}
