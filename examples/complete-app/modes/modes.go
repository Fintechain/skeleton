// Package modes contains different execution modes for the complete application example.
// This version uses the new Builder API instead of FX dependency injection.
package modes

import (
	"fmt"

	"github.com/fintechain/skeleton/examples/complete-app/plugins/database"
	"github.com/fintechain/skeleton/examples/complete-app/plugins/webserver"
	"github.com/fintechain/skeleton/internal/domain/config"
	"github.com/fintechain/skeleton/internal/domain/event"
	"github.com/fintechain/skeleton/internal/domain/logging"
	infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
	infraEvent "github.com/fintechain/skeleton/internal/infrastructure/event"
	infraLogging "github.com/fintechain/skeleton/internal/infrastructure/logging"
	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/runtime"
)

// RunDaemonMode demonstrates running multiple plugins as long-running services using Builder API
func RunDaemonMode() error {
	fmt.Println("=== Complete Application - Daemon Mode (Builder API) ===")
	fmt.Println("Starting multiple plugins with services using Builder pattern...")

	// Start daemon with multiple plugins using Builder API
	return runtime.NewBuilder().
		WithPlugins(
			webserver.NewWebServerPlugin(8080),
			database.NewDatabasePlugin("postgres", "test://connection"),
		).
		BuildDaemon()
}

// RunCommandMode demonstrates executing one operation and exiting using Builder API
func RunCommandMode() error {
	fmt.Println("=== Complete Application - Command Mode (Builder API) ===")

	// Execute one operation and exit using Builder API
	result, err := runtime.NewBuilder().
		WithPlugins(
			webserver.NewWebServerPlugin(8080),
			database.NewDatabasePlugin("postgres", "test://connection"),
		).
		BuildCommand("database-query", map[string]interface{}{
			"query": "SELECT * FROM users WHERE active = true",
		})

	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
		return err
	}

	fmt.Printf("Operation result: %+v\n", result)
	fmt.Println("Command mode complete - application exits")
	return nil
}

// RunWithCustomProviders demonstrates using custom implementations of framework services using Builder API
func RunWithCustomProviders() error {
	fmt.Println("=== Complete Application - Custom Dependencies (Builder API) ===")
	fmt.Println("Using custom configuration, logger, and event bus with Builder pattern...")

	// Create custom dependencies directly (no FX complexity!)
	customConfig := createCustomConfiguration()
	customLogger := createCustomLogger()
	customEventBus := createCustomEventBus()

	fmt.Println("🔍 Testing if custom configuration is working...")
	fmt.Printf("   Custom app.name: '%s'\n", customConfig.GetStringDefault("app.name", "NOT FOUND"))
	fmt.Printf("   Custom app.version: '%s'\n", customConfig.GetStringDefault("app.version", "NOT FOUND"))
	fmt.Printf("   Custom database.host: '%s'\n", customConfig.GetStringDefault("database.host", "NOT FOUND"))

	// Test all expected values
	fmt.Println("🔍 Testing all custom configuration values...")
	fmt.Printf("   app.name: '%s'\n", customConfig.GetString("app.name"))
	fmt.Printf("   app.version: '%s'\n", customConfig.GetString("app.version"))
	fmt.Printf("   database.host: '%s'\n", customConfig.GetString("database.host"))
	fmt.Printf("   database.port: %d\n", customConfig.GetIntDefault("database.port", 0))
	fmt.Printf("   builder.api: %t\n", customConfig.GetBoolDefault("builder.api", false))

	// Check if configuration exists
	fmt.Printf("   app.name exists: %t\n", customConfig.Exists("app.name"))
	fmt.Printf("   app.version exists: %t\n", customConfig.Exists("app.version"))

	// Execute with custom dependencies using simple Builder API
	result, err := runtime.NewBuilder().
		WithConfig(customConfig).
		WithLogger(customLogger).
		WithEventBus(customEventBus).
		WithPlugins(
			webserver.NewWebServerPlugin(8080),
			database.NewDatabasePlugin("postgres", "custom://connection"),
		).
		BuildCommand("database-query", map[string]interface{}{
			"query": "SELECT * FROM users WHERE active = true",
		})

	if err != nil {
		fmt.Printf("Operation with custom dependencies failed: %v\n", err)
		return err
	}

	fmt.Printf("Operation with custom dependencies result: %+v\n", result)

	// Verify if custom config values appear in the result
	if appName, ok := result["app_name"].(string); ok {
		if appName == "CUSTOM-BUILDER-API-APP" {
			fmt.Println("✅ SUCCESS: Custom configuration IS working! Found custom app name.")
		} else {
			fmt.Printf("❌ ISSUE: Custom configuration might not be working. Expected 'CUSTOM-BUILDER-API-APP', got '%s'\n", appName)
		}
	}

	fmt.Println("Custom dependencies example complete - Builder API made it simple!")
	return nil
}

// createCustomConfiguration creates a custom configuration implementation
// This replaces the complex FX provider pattern with simple direct creation
func createCustomConfiguration() config.Configuration {
	fmt.Println("🔧 Creating CUSTOM CONFIGURATION with unique values...")

	data := map[string]interface{}{
		"app.name":       "CUSTOM-BUILDER-API-APP", // Very distinctive name
		"app.version":    "99.99.99",               // Very distinctive version
		"database.host":  "custom-builder-db-host",
		"database.port":  9999,
		"server.port":    8080,
		"logging.level":  "debug",
		"features.cache": true,
		"builder.api":    true, // New feature flag showing Builder API usage
	}

	fmt.Printf("🔍 Input data: %+v\n", data)

	config := infraConfig.NewMemoryConfigurationWithData(data)

	// Test the config immediately after creation
	fmt.Println("🔍 Testing config immediately after creation:")
	fmt.Printf("   config.GetString(\"app.name\"): \"%s\"\n", config.GetString("app.name"))
	fmt.Printf("   config.Exists(\"app.name\"): %t\n", config.Exists("app.name"))

	// Test direct source access
	source := config.GetSource()
	value, exists := source.GetValue("app.name")
	fmt.Printf("   source.GetValue(\"app.name\"): value=%v, exists=%t\n", value, exists)

	return config
}

// createCustomLogger creates a custom logger implementation
// This replaces the complex FX provider pattern with simple direct creation
func createCustomLogger() logging.LoggerService {
	fmt.Println("🔧 Creating CUSTOM LOGGER...")

	config := component.ComponentConfig{
		ID:          "custom-logger-complete-app",
		Name:        "Custom Logger for Complete App",
		Description: "Custom logger demonstrating Builder API simplicity",
		Type:        component.TypeService,
		Version:     "1.0.0",
	}

	// Create a NoOp logger wrapped in the service (could be any logger implementation)
	noOpLogger := infraLogging.NewNoOpLogger()
	customLogger, err := infraLogging.NewLogger(config, noOpLogger)
	if err != nil {
		// In a real application, you'd handle this more gracefully
		panic(fmt.Sprintf("Failed to create custom logger: %v", err))
	}

	// Test the logger immediately
	fmt.Printf("🔍 Custom logger ID: %s\n", customLogger.ID())
	fmt.Printf("🔍 Custom logger Name: %s\n", customLogger.Name())

	return customLogger
}

// createCustomEventBus creates a custom event bus implementation
// This replaces the complex FX provider pattern with simple direct creation
func createCustomEventBus() event.EventBusService {
	fmt.Println("🔧 Creating CUSTOM EVENT BUS...")

	config := component.ComponentConfig{
		ID:          "custom-eventbus-complete-app",
		Name:        "Custom Event Bus for Complete App",
		Description: "Custom event bus demonstrating Builder API simplicity",
		Type:        component.TypeService,
		Version:     "1.0.0",
	}

	eventBus := infraEvent.NewEventBus(config)

	// Test the event bus immediately
	fmt.Printf("🔍 Custom event bus ID: %s\n", eventBus.ID())
	fmt.Printf("🔍 Custom event bus Name: %s\n", eventBus.Name())

	return eventBus
}
