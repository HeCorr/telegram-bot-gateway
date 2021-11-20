package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/apparentlymart/go-cidr/cidr"
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
