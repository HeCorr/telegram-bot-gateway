package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

func main() {
	botsFile := flag.String("f", findBotsFile(), "Use the specified .yaml file")
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

	fmt.Println("Loaded routes:")
	e := echo.New()
	e.HideBanner = true


	for _, b := range botsData.Bots {
		fmt.Printf("    %"+strconv.Itoa(namePd)+"s: %-"+strconv.Itoa(endpointPd)+"s -> %s\n", b.Name, b.Endpoint, b.Path)
	}

	if err = e.Start(":9000"); err != nil {
		fmt.Println(err)
	}
}

func registerRoute(e *echo.Echo, endpoint string, path string) {
	e.POST(endpoint, func(c echo.Context) error {
		resp, err := http.Post(path, c.Request().Header.Get("Content-Type"), c.Request().Body)
		if err != nil {
			// TODO: improve logging
			fmt.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(resp.StatusCode)
	})
}
