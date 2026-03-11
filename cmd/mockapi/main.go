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

	// Simple CLI: mockapi [port]
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help", "--help", "-h":
			fmt.Println("MockAPI — Lightweight local API mock server")
			fmt.Println("")
			fmt.Println("Usage:")
			fmt.Println("  mockapi [port]        Start server (default: 8088)")
			fmt.Println("  mockapi --port 3000   Start server on custom port")
			fmt.Println("")
			fmt.Println("Routes are saved to mockapi.json")
			fmt.Println("Web UI at http://localhost:<port>/")
			fmt.Println("Mock endpoints at http://localhost:<port>/mock/<path>")
			return
		case "--port":
			if len(os.Args) < 3 {
				fmt.Println("Usage: mockapi --port <number>")
				os.Exit(1)
			}
			fmt.Sscanf(os.Args[2], "%d", &port)
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

	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
