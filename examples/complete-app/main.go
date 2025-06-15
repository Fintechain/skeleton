// Package main demonstrates a complete application using the Fintechain Skeleton framework
// showcasing plugins, custom providers, and all framework patterns.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fintechain/skeleton/examples/complete-app/modes"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Complete Framework Example - Fintechain Skeleton Framework")
		fmt.Println()
		fmt.Println("This example demonstrates:")
		fmt.Println("  • Plugin orchestration (webserver + database)")
		fmt.Println("  • Custom providers (logger, config, eventbus)")
		fmt.Println("  • Component lifecycle management")
		fmt.Println("  • Both daemon and command modes")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  go run examples/complete-app/main.go daemon     # Run as daemon (long-running services)")
		fmt.Println("  go run examples/complete-app/main.go command    # Run in command mode (execute and exit)")
		fmt.Println("  go run examples/complete-app/main.go custom     # Run with custom providers")
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "daemon":
		if err := modes.RunDaemonMode(); err != nil {
			log.Fatalf("Daemon mode failed: %v", err)
		}
	case "command":
		if err := modes.RunCommandMode(); err != nil {
			log.Fatalf("Command mode failed: %v", err)
		}
	case "custom":
		if err := modes.RunWithCustomProviders(); err != nil {
			log.Fatalf("Custom providers mode failed: %v", err)
		}
	default:
		fmt.Printf("Unknown mode: %s\n", mode)
		os.Exit(1)
	}
}
