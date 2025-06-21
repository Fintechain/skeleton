// Package main demonstrates a complete application using the Fintechain Skeleton framework
// showcasing plugins, custom dependencies, and all framework patterns using the new Builder API.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fintechain/skeleton/examples/complete-app/modes"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Complete Framework Example - Fintechain Skeleton Framework (Builder API)")
		fmt.Println()
		fmt.Println("This example demonstrates the new Builder API which replaces FX dependency injection:")
		fmt.Println("  • Simple, explicit dependency injection (no FX complexity)")
		fmt.Println("  • Plugin orchestration (webserver + database)")
		fmt.Println("  • Custom dependencies (logger, config, eventbus)")
		fmt.Println("  • Component lifecycle management")
		fmt.Println("  • Both daemon and command modes")
		fmt.Println()
		fmt.Println("Advantages of Builder API over FX:")
		fmt.Println("  ✅ No FX knowledge required")
		fmt.Println("  ✅ Easy debugging and testing")
		fmt.Println("  ✅ Clear dependency flow")
		fmt.Println("  ✅ Compile-time errors for missing dependencies")
		fmt.Println("  ✅ Simple custom dependency injection")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  go run examples/complete-app/main.go daemon     # Run as daemon (long-running services)")
		fmt.Println("  go run examples/complete-app/main.go command    # Run in command mode (execute and exit)")
		fmt.Println("  go run examples/complete-app/main.go custom     # Run with custom dependencies")
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
			log.Fatalf("Custom dependencies mode failed: %v", err)
		}
	default:
		fmt.Printf("Unknown mode: %s\n", mode)
		fmt.Println("Use: daemon, command, or custom")
		os.Exit(1)
	}
}
