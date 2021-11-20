package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var telegramCIDRs = []string{"149.154.160.0/20", "91.108.4.0/22"}

func telegramIPMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		for _, cidr := range telegramCIDRs {
			inRange, err := ipInCIDR(c.RealIP(), cidr)
			if err != nil {
				fmt.Printf("Failed to check IP %s: %v\n", c.RealIP(), err)
				return c.NoContent(http.StatusInternalServerError)
			}
			if inRange {
				return next(c)
			}
		}
		return c.NoContent(http.StatusForbidden)
	}
}
