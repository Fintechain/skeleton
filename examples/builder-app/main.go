// Package main demonstrates the new Builder API for the Fintechain Skeleton framework.
//
// This example shows how to use the builder pattern instead of FX dependency injection
// for creating and running Fintechain applications.
//
// Usage:
//
//	go run main.go daemon    # Run in daemon mode
//	go run main.go command   # Run in command mode
//	go run main.go custom    # Run custom dependencies examples
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fintechain/skeleton/pkg/runtime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [daemon|command|custom]")
		fmt.Println("  daemon  - Run as long-running service")
		fmt.Println("  command - Execute operation and exit")
		fmt.Println("  custom  - Demonstrate custom dependency injection")
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "daemon":
		runDaemonMode()
	case "command":
		runCommandMode()
	case "custom":
		runCustomExamples()
	default:
		fmt.Printf("Unknown mode: %s\n", mode)
		fmt.Println("Use 'daemon', 'command', or 'custom'")
		os.Exit(1)
	}
}

// runDaemonMode demonstrates running a daemon application using the builder API
func runDaemonMode() {
	fmt.Println("=== Builder API Daemon Mode Example ===")

	// Create a simple test plugin
	testPlugin := NewTestPlugin()

	// Use the builder API to create and run a daemon
	err := runtime.NewBuilder().
		WithPlugins(testPlugin).
		BuildDaemon()

	if err != nil {
		log.Fatal("Failed to start daemon:", err)
	}
}

// runCommandMode demonstrates running a command application using the builder API
func runCommandMode() {
	fmt.Println("=== Builder API Command Mode Example ===")

	// Create a simple test plugin
	testPlugin := NewTestPlugin()

	// Use the builder API to execute a command
	result, err := runtime.NewBuilder().
		WithPlugins(testPlugin).
		BuildCommand("test-operation", map[string]interface{}{
			"message": "Hello from Builder API!",
			"number":  42,
		})

	if err != nil {
		log.Fatal("Failed to execute command:", err)
	}

	fmt.Printf("Command result: %+v\n", result)
}
