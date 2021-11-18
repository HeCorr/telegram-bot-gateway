package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Reads the file f and parses the YAML data into bots
func readBotsFile(f string) (bots Bots, _ error) {
	data, err := os.ReadFile(f)
	if err != nil {
		return bots, err
	}
	err = yaml.Unmarshal(data, &bots)
	return bots, err
}
