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

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help", "--help", "-h":
			fmt.Println("MockAPI — Lightweight local API mock server")
			fmt.Println("")
			fmt.Println("Usage:")
			fmt.Println("  mockapi [port]          Start server (default: 8088)")
			fmt.Println("  mockapi --port 3000     Start on custom port")
			fmt.Println("  mockapi --config file   Use custom config file")
			fmt.Println("")
			fmt.Println("Features:")
			fmt.Println("  • Web UI for visual route management")
			fmt.Println("  • Dynamic path params (/users/:id) with {{id}} substitution")
			fmt.Println("  • Wildcard paths (/api/*)")
			fmt.Println("  • Request logging with auto-refresh")
			fmt.Println("  • Import/Export routes as JSON")
			fmt.Println("  • Preset templates for common APIs")
			fmt.Println("  • CORS enabled by default")
			fmt.Println("  • Simulated response delay")
			return
		case "--port":
			if len(os.Args) < 3 {
				fmt.Println("Usage: mockapi --port <number>")
				os.Exit(1)
			}
			fmt.Sscanf(os.Args[2], "%d", &port)
		case "--config":
			if len(os.Args) < 3 {
				fmt.Println("Usage: mockapi --config <file>")
				os.Exit(1)
			}
			configFile = os.Args[2]
		default:
			fmt.Sscanf(os.Args[1], "%d", &port)
		}
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	cfg.Port = port

	srv := server.New(cfg, configFile)
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
