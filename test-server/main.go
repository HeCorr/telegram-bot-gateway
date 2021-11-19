package main

import (
	"flag"
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	port := flag.String("p", "9800", "Listen port")
	flag.Parse()

	e := echo.New()

	e.POST("/bot", func(c echo.Context) error {
		var asd TelegramUpdate
		err := c.Bind(&asd)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("%s: %s\n", asd.Message.Chat.Username, asd.Message.Text)
		return nil
	})

	fmt.Println(e.Start(":" + *port))
}
