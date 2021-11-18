package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	botsFile := flag.String("f", findBotsFile(), "Use the specified .yaml file")
	flag.Parse()

	if *botsFile == "" {
		fmt.Println("No bots.yml file found. Please create one.")
		os.Exit(1)
	}
}
