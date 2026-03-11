package main

import (
	"fmt"
	"os"

	"mockapi/pkg/config"
	"mockapi/pkg/server"
)

func main() {
	port := 8088
	configFile := "mockapi.json"
	useHTTPS := false
	certFile := ""
	keyFile := ""

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "help", "--help", "-h":
			fmt.Println("MockAPI — Lightweight local API mock server")
			fmt.Println("")
			fmt.Println("Usage:")
			fmt.Println("  mockapi [port]              Start server (default: 8088)")
			fmt.Println("  mockapi --port 3000         Start on custom port")
			fmt.Println("  mockapi --config file.json  Use custom config file")
			fmt.Println("  mockapi --https --cert c.pem --key k.pem  Enable HTTPS")
			fmt.Println("  mockapi --proxy http://backend  Set proxy backend")
			fmt.Println("")
			fmt.Println("Features:")
			fmt.Println("  • Web UI for visual route management")
			fmt.Println("  • Dynamic path params (/users/:id) with {{id}} substitution")
			fmt.Println("  • Wildcard paths (/api/*), ALL methods")
			fmt.Println("  • Conditional responses (match headers/body)")
			fmt.Println("  • Proxy mode (forward unmatched to real backend)")
			fmt.Println("  • Request logging with auto-refresh")
			fmt.Println("  • Import/Export routes as JSON")
			fmt.Println("  • Preset templates for common APIs")
			fmt.Println("  • HTTPS support")
			fmt.Println("  • CORS enabled by default")
			fmt.Println("  • Simulated response delay")
			return
		case "--port":
			i++
			if i >= len(args) {
				fmt.Println("Usage: mockapi --port <number>")
				os.Exit(1)
			}
			fmt.Sscanf(args[i], "%d", &port)
		case "--config":
			i++
			if i >= len(args) {
				fmt.Println("Usage: mockapi --config <file>")
				os.Exit(1)
			}
			configFile = args[i]
		case "--https":
			useHTTPS = true
		case "--cert":
			i++
			certFile = args[i]
		case "--key":
			i++
			keyFile = args[i]
		case "--proxy":
			i++
			// Save proxy URL directly
			cfg, _ := config.Load(configFile)
			if cfg != nil {
				cfg.ProxyURL = args[i]
				cfg.Save(configFile)
			}
		default:
			fmt.Sscanf(args[i], "%d", &port)
		}
	}

	if useHTTPS && (certFile == "" || keyFile == "") {
		fmt.Println("HTTPS requires --cert and --key flags")
		os.Exit(1)
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	cfg.Port = port

	srv := server.New(cfg, configFile)

	if useHTTPS {
		srv.EnableHTTPS(certFile, keyFile)
	}

	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
