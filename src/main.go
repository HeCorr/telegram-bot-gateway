package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var client *http.Client

func init() {
	// clone default transport to keep default settings
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{
		// don't verify certificates
		InsecureSkipVerify: true,
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5,
	}
}

func main() {
	botsFile := flag.String("f", findBotsFile(), "Use the specified .yaml file")
	listenAddr := flag.String("l", "", "Listen address (default \"localhost:8443\")")
	certFile := flag.String("c", "", "Certificate file for HTTPS")
	keyFile := flag.String("k", "", "Private key file for HTTPS")
	strict := flag.Bool("s", false, "Strict mode - blocks requests not coming from Telegram")
	init := flag.Bool("i", false, "Initialize (create) example bots.yaml file")

	flag.Parse()

	if *init {
		if err := initBotsFile(); err != nil {
			fmt.Printf("Failed to initialize bots.yaml file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Example bots.yaml file successfuly created! Please edit it and run again.")
		os.Exit(0)
	}

	if *botsFile == "" {
		fmt.Println("Default bots.yaml file not found. Please create it with -i or specify one with -f.")
		os.Exit(1)
	}

	botsData, err := readBotsFile(*botsFile)
	if err != nil {
		fmt.Printf("Failed to load bots file: %v", err)
		os.Exit(1)
	}

	if *keyFile == "" {
		if botsData.Options.Key == "" {
			fmt.Println("Private key file not specified. Please specify it with the -k flag or in the bots.yaml file.")
			os.Exit(1)
		}
		*keyFile = botsData.Options.Key
	}

	if *certFile == "" {
		if botsData.Options.Cert == "" {
			fmt.Println("Certificate file not specified. Please specify it with the -c flag or in the bots.yaml file")
			os.Exit(1)
		}
		*certFile = botsData.Options.Cert
	}

	if *listenAddr == "" {
		if botsData.Options.Listen == "" {
			*listenAddr = "localhost:8443"
		} else {
			*listenAddr = botsData.Options.Listen
		}
	}

	// calculate dynamic padding based on the longest values
	var namePd, endpointPd int
	for _, b := range botsData.Bots {
		if namePd < len(b.Name) {
			namePd = len(b.Name)
		}
		if endpointPd < len(b.Endpoint) {
			endpointPd = len(b.Endpoint)
		}
	}

	e := echo.New()
	e.HideBanner = true

	// https://echo.labstack.com/guide/ip-address/
	// TODO: add setting for controlling this
	e.IPExtractor = echo.ExtractIPDirect()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "[${time_custom}] ${remote_ip} -> ${path} [${status}] (${latency_human})\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	if *strict || botsData.Options.Strict {
		fmt.Println("Strict mode enabled!")
		e.Use(telegramIPMiddleware)
	}

	fmt.Println("Loading routes...")

	for _, b := range botsData.Bots {
		b.NormalizeEndpoint()
		fmt.Printf("    %"+strconv.Itoa(namePd)+"s: %-"+strconv.Itoa(endpointPd)+"s -> %s\n", b.Name, b.Endpoint, b.Path)
		registerRoute(e, b.Endpoint, b.Path)
	}

	go func() {
		err := e.StartTLS(*listenAddr, *certFile, *keyFile)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server start failed:", err)
			os.Exit(1)
		}
	}()

	// Ctrl + C handling
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)

	<-exitChan
	fmt.Println("Interrupt received, stopping webserver...")
	// TODO: implement pressing Ctrl + C again for forcing exit
	if err = e.Shutdown(context.Background()); err != nil {
		fmt.Println("Server shutdown failed:", err)
	}
}
