package main

import (
	"fmt"
	"os"
)

func main() {
	botsFile := findBotsFile()
	if botsFile == "" {
		fmt.Println("No bots.yml file found. Please create one.")
		os.Exit(1)
	}
}
