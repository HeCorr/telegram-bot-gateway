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
		fmt.Println("Default bots.yml file not found. Please create it or specify one with -f.")
		os.Exit(1)
	}
}
