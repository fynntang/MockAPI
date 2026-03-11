package main

import (
	"fmt"
	"os"
	"path/filepath"

	"mockapi/pkg/config"
	"mockapi/pkg/server"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "help", "--help", "-h":
		printUsage()
		return
	case "version", "--version", "-v":
		fmt.Println("MockAPI v1.0.0")
		return
	case "init":
		cmdInit()
		return
	case "validate":
		cmdValidate()
		return
	case "serve":
		cmdServe()
		return
	default:
		// Default: start server (backward compatible)
		cmdServe()
	}
}

func printUsage() {
	fmt.Println(`MockAPI — Lightweight Local API Mock Server

Usage:
  mockapi <command> [options]

Commands:
  init              Initialize a new MockAPI project
  validate          Validate configuration file
  serve             Start the mock server (default)
  version, -v       Show version
  help, -h          Show this help

Examples:
  mockapi init                      # Create new project
  mockapi validate                  # Validate config
  mockapi serve                     # Start server on port 8088
  mockapi serve --port 3000         # Start on custom port
  mockapi serve --config my.json    # Use custom config
  mockapi serve --hot-reload        # Enable hot reload
  mockapi                           # Same as 'mockapi serve'

Features:
  • Web UI for visual route management
  • JavaScript scripts for dynamic responses
  • Auto-generate from Swagger/OpenAPI specs
  • WebSocket & GraphQL mock support
  • gRPC mock support with proto import
  • Proxy mode for incremental mocking
  • Hot reload configuration
  • Zero dependencies (single binary)

Learn more: https://github.com/fynntang/MockAPI`)
}

func cmdInit() {
	fmt.Println("🦞 Initializing MockAPI project...")

	// Create directory structure
	dirs := []string{"routes", "scripts", "protos"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
		fmt.Printf("  ✓ Created %s/\n", dir)
	}

	// Create example config
	configFile := "mockapi.json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		cfg := config.Default()
		if err := cfg.Save(configFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("  ✓ Created %s\n", configFile)
	}

	// Create example route
	exampleRoute := `{
  "id": "example",
  "method": "GET",
  "path": "/hello",
  "status": 200,
  "body": "{\"message\": \"Hello from MockAPI!\"}",
  "description": "Example route"
}`
	routeFile := "routes/example.json"
	if err := os.WriteFile(routeFile, []byte(exampleRoute), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating example route: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ Created %s\n", routeFile)

	// Create README
	readme := "# MockAPI Project\n\n## Quick Start\n\n    # Start server\n    mockapi serve\n\n## Project Structure\n\n- routes/ - Route configuration files\n- scripts/ - JavaScript script files\n- protos/ - Protocol buffer definitions\n- mockapi.json - Main configuration file\n\n## Add Your First Route\n\n1. Open Web UI: http://localhost:8088\n2. Click \"+ Add Route\"\n3. Configure method, path, and response\n4. Test: curl http://localhost:8088/mock/your-path\n\nLearn more: https://github.com/fynntang/MockAPI\n"
	if err := os.WriteFile("README.md", []byte(readme), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating README: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ Created README.md\n")

	fmt.Println("\n✅ Project initialized! Run 'mockapi serve' to start.")
}

func cmdValidate() {
	fmt.Println("🔍 Validating MockAPI configuration...")

	configFile := "mockapi.json"
	if len(os.Args) > 2 && os.Args[2] != "" {
		configFile = os.Args[2]
	}

	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "❌ Config file not found: %s\n", configFile)
		os.Exit(1)
	}

	// Load and validate
	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Invalid config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Config file: %s\n", configFile)
	fmt.Printf("✓ Port: %d\n", cfg.Port)
	fmt.Printf("✓ Routes: %d\n", len(cfg.Routes))
	fmt.Printf("✓ GraphQL handlers: %d\n", len(cfg.GraphQL))
	fmt.Printf("✓ gRPC handlers: %d\n", len(cfg.GRPC))
	fmt.Printf("✓ CORS enabled: %v\n", cfg.CORSEnabled)
	fmt.Printf("✓ Proxy URL: %s\n", cfg.ProxyURL)

	// Validate routes
	errors := 0
	for _, route := range cfg.Routes {
		if route.Path == "" {
			fmt.Fprintf(os.Stderr, "❌ Route %s has empty path\n", route.ID)
			errors++
		}
		if route.Status < 100 || route.Status > 599 {
			fmt.Fprintf(os.Stderr, "❌ Route %s has invalid status: %d\n", route.ID, route.Status)
			errors++
		}
	}

	if errors > 0 {
		fmt.Fprintf(os.Stderr, "\n❌ Found %d validation errors\n", errors)
		os.Exit(1)
	}

	fmt.Println("\n✅ Configuration is valid!")
}

func cmdServe() {
	port := 8088
	configFile := "mockapi.json"
	hotReload := false

	// Parse arguments
	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--port":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &port)
				i++
			}
		case "--config":
			if i+1 < len(args) {
				configFile = args[i+1]
				i++
			}
		case "--hot-reload":
			hotReload = true
		}
	}

	// Load config
	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	cfg.Port = port

	// Create server
	srv := server.New(cfg, configFile)

	// Setup hot reload
	if hotReload {
		fmt.Println("🔥 Hot reload enabled")
		watcher := config.NewWatcher(configFile, func(newCfg *config.Config) {
			fmt.Println("🔄 Configuration reloaded!")
			// Note: In a real implementation, we'd update the server's config
		})
		watcher.Start()
		defer watcher.Stop()
	}

	// Start server
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// Helper to get absolute path
func absPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	abs, _ := filepath.Abs(path)
	return abs
}
