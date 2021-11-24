package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"
)

var (
	defaultBotsFiles = []string{"bots.yaml", "bots.yml"}
)

// Checks if file f exists and is not a directory
// (it panicks if the error is not ErrNotExist)
func fileExists(f string) bool {
	info, err := os.Stat(f)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		} else {
			panic(err)
		}
	}
	return !info.IsDir()
}

// Looks for the default bots file (in the current path) in the following order:
// - bots.yaml
// - bots.yml
func findBotsFile() string {
	for _, f := range defaultBotsFiles {
		if fileExists(f) {
			return f
		}
	}
	return ""
}

// Reads the file f and parses the YAML data into bots
func readBotsFile(f string) (bots Bots, _ error) {
	data, err := os.ReadFile(f)
	if err != nil {
		return bots, err
	}
	err = yaml.Unmarshal(data, &bots)
	if err != nil {
		return bots, err
	}
	for i := 0; i < len(bots.Bots); i++ {
		if bots.Bots[i].Name == "" {
			bots.Bots[i].Name = "Bot " + strconv.Itoa(i+1)
		}
	}
	return bots, nil
}

// Checks if CIDR contains IP
func ipInCIDR(IP, CIDR string) (bool, error) {
	ip := net.ParseIP(IP).To4()
	if ip == nil {
		return false, fmt.Errorf("parse %s: invalid ipv4", IP)
	}
	_, network, err := net.ParseCIDR(CIDR)
	if err != nil {
		return false, err
	}
	start, end := cidr.AddressRange(network)
	if bytes.Compare(ip, start) >= 0 && bytes.Compare(ip, end) <= 0 {
		return true, nil
	}
	return false, nil
}

// Register route which forwards requests from endpoint to path
func registerRoute(e *echo.Echo, endpoint string, path string) {
	e.POST(endpoint, func(c echo.Context) error {
		resp, err := client.Post(path, c.Request().Header.Get("Content-Type"), c.Request().Body)
		if err != nil {
			// TODO: improve logging
			fmt.Println(err)
			return c.NoContent(http.StatusGatewayTimeout)
		}
		defer resp.Body.Close()
		return c.NoContent(resp.StatusCode)
	})
}
