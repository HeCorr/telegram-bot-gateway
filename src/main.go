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
	botsData, err := readBotsFile(*botsFile)
	if err != nil {
		fmt.Printf("Failed to load bots file: %v", err)
		os.Exit(1)
	}
	fmt.Println("Loaded routes:")
	for _, b := range botsData.Bots {
		// TODO: calculate padding based on longest value
		fmt.Printf("  %s: %s -> %s\n", b.Name, b.Endpoint, b.Path)
	}
}
