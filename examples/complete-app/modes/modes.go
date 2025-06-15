// Package modes contains different execution modes for the complete application example.
package modes

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/fintechain/skeleton/examples/complete-app/plugins/database"
	"github.com/fintechain/skeleton/examples/complete-app/plugins/webserver"
	"github.com/fintechain/skeleton/examples/complete-app/providers"
	"github.com/fintechain/skeleton/pkg/runtime"
)

// RunDaemonMode demonstrates running multiple plugins as long-running services
func RunDaemonMode() error {
	fmt.Println("=== Complete Application - Daemon Mode ===")
	fmt.Println("Starting multiple plugins with services...")

	// Start daemon with multiple plugins using default providers
	return runtime.StartDaemon(
		runtime.WithPlugins(
			webserver.NewWebServerPlugin(8080),
			database.NewDatabasePlugin("postgres", "test://connection")),
	)
}

// RunCommandMode demonstrates executing one operation and exiting
func RunCommandMode() error {
	fmt.Println("=== Complete Application - Command Mode ===")

	// Execute one operation and exit
	result, err := runtime.ExecuteCommand("http-route", map[string]any{
		"method": "GET",
		"path":   "/api/health",
	}, runtime.WithPlugins(
		webserver.NewWebServerPlugin(8080),
		database.NewDatabasePlugin("postgres", "test://connection")),
	)

	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
		return err
	}

	fmt.Printf("Operation result: %+v\n", result)
	fmt.Println("Command mode complete - application exits")
	return nil
}

// RunWithCustomProviders demonstrates using custom implementations of framework services
func RunWithCustomProviders() error {
	fmt.Println("=== Complete Application - Custom Providers ===")
	fmt.Println("Using custom logger, configuration, and event bus...")

	// Execute with custom providers to show the full power of FX integration
	result, err := runtime.ExecuteCommand("database-query", map[string]any{
		"query": "SELECT * FROM users WHERE active = true",
	},
		// Plugins
		runtime.WithPlugins(
			webserver.NewWebServerPlugin(8080),
			database.NewDatabasePlugin("postgres", "custom://connection")),

		// Custom providers using FX options
		runtime.WithOptions(
			fx.Replace(providers.NewCustomLogger()),
			fx.Replace(providers.NewCustomConfiguration()),
			fx.Replace(providers.NewCustomEventBus()),
		),
	)

	if err != nil {
		fmt.Printf("Operation with custom providers failed: %v\n", err)
		return err
	}

	fmt.Printf("Operation with custom providers result: %+v\n", result)
	fmt.Println("Custom providers example complete")
	return nil
}
