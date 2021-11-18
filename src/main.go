package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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

	for _, b := range botsData.Bots {
		fmt.Printf("  %"+strconv.Itoa(namePd)+"s: %-"+strconv.Itoa(endpointPd)+"s -> %s\n", b.Name, b.Endpoint, b.Path)
	}
}
