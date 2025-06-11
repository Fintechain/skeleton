// Package main demonstrates a complete application using the Fintechain Skeleton framework
// with multiple plugins working together.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fintechain/skeleton/examples/plugins/database"
	"github.com/fintechain/skeleton/examples/plugins/webserver"
	"github.com/fintechain/skeleton/pkg/fx"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Complete Application Example - Fintechain Skeleton Framework")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  go run examples/complete-app/main.go daemon     # Run as daemon (long-running services)")
		fmt.Println("  go run examples/complete-app/main.go command    # Run in command mode (execute and exit)")
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "daemon":
		runDaemonMode()
	case "command":
		runCommandMode()
	default:
		fmt.Printf("Unknown mode: %s\n", mode)
		os.Exit(1)
	}
}

// runDaemonMode demonstrates running multiple plugins as long-running services
func runDaemonMode() {
	fmt.Println("=== Complete Application - Daemon Mode ===")
	fmt.Println("Starting multiple plugins with services...")

	// Start daemon with multiple plugins
	// This demonstrates how multiple plugins work together in daemon mode
	err := fx.StartDaemon(
		fx.WithPlugins(
			webserver.NewWebServerPlugin(8080),
			database.NewDatabasePlugin("postgres", "test://connection")),
	)

	if err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}

	// Framework handles:
	// 1. Plugin initialization and component registration
	// 2. Service startup (HTTPService and DatabaseConnectionService)
	// 3. Services run continuously until shutdown signal
	// 4. Graceful shutdown on SIGINT/SIGTERM
}

// runCommandMode demonstrates executing one operation and exiting
func runCommandMode() {
	fmt.Println("=== Complete Application - Command Mode ===")

	// Execute one operation and exit
	result, err := fx.ExecuteCommand("http-route", map[string]any{
		"method": "GET",
		"path":   "/api/health",
	}, fx.WithPlugins(
		webserver.NewWebServerPlugin(8080),
		database.NewDatabasePlugin("postgres", "test://connection")),
	)

	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
	} else {
		fmt.Printf("Operation result: %+v\n", result)
	}

	fmt.Println("Command mode complete - application exits")
}
