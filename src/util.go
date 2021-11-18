package main

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	defaultBotsFiles = []string{"bots.yml", "bots.yaml"}
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
// - bots.yml
// - bots.yaml
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
	return bots, err
}
