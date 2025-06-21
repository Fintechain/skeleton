// Package main provides examples of custom dependency injection with the Builder API.
package main

import (
	"fmt"
	"log"

	"github.com/fintechain/skeleton/internal/domain/logging"
	infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
	infraLogging "github.com/fintechain/skeleton/internal/infrastructure/logging"
	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/runtime"
)

// demonstrateCustomConfig shows how to inject a custom configuration
func demonstrateCustomConfig() {
	fmt.Println("=== Custom Configuration Example ===")

	// Create custom configuration with application-specific settings
	customConfig := infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
		"app.name":               "Custom Builder App",
		"app.version":            "2.0.0",
		"test_service.enabled":   true,
		"test_service.interval":  "10s",
		"custom.feature_flag":    true,
		"custom.max_connections": 100,
	})

	// Create test plugin
	testPlugin := NewTestPlugin()

	// Use builder with custom configuration
	result, err := runtime.NewBuilder().
		WithConfig(customConfig).
		WithPlugins(testPlugin).
		BuildCommand("test-operation", map[string]interface{}{
			"message": "Testing custom config",
			"number":  123,
		})

	if err != nil {
		log.Fatal("Failed to execute command with custom config:", err)
	}

	fmt.Printf("Result with custom config: %+v\n", result)
	fmt.Printf("App name from config: %s\n", customConfig.GetString("app.name"))
	fmt.Printf("App version from config: %s\n", customConfig.GetString("app.version"))
}

// demonstrateCustomLogger shows how to inject a custom logger
func demonstrateCustomLogger() {
	fmt.Println("\n=== Custom Logger Example ===")

	// Create a custom logger (using NoOp for simplicity, but could be any implementation)
	customLogger := createCustomLogger()

	// Create test plugin
	testPlugin := NewTestPlugin()

	// Use builder with custom logger
	result, err := runtime.NewBuilder().
		WithLogger(customLogger).
		WithPlugins(testPlugin).
		BuildCommand("test-operation", map[string]interface{}{
			"message": "Testing custom logger",
			"number":  456,
		})

	if err != nil {
		log.Fatal("Failed to execute command with custom logger:", err)
	}

	fmt.Printf("Result with custom logger: %+v\n", result)
	fmt.Println("Custom logger was used for all framework logging")
}

// demonstrateFullCustomization shows how to inject multiple custom dependencies
func demonstrateFullCustomization() {
	fmt.Println("\n=== Full Custom Dependencies Example ===")

	// Create custom configuration
	customConfig := infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
		"app.name":          "Fully Customized App",
		"app.environment":   "production",
		"logging.level":     "info",
		"features.advanced": true,
	})

	// Create custom logger
	customLogger := createCustomLogger()

	// Create test plugin
	testPlugin := NewTestPlugin()

	// Use builder with multiple custom dependencies
	result, err := runtime.NewBuilder().
		WithConfig(customConfig).
		WithLogger(customLogger).
		WithPlugins(testPlugin).
		BuildCommand("test-operation", map[string]interface{}{
			"message": "Testing full customization",
			"number":  789,
		})

	if err != nil {
		log.Fatal("Failed to execute command with full customization:", err)
	}

	fmt.Printf("Result with full customization: %+v\n", result)
	fmt.Printf("Environment: %s\n", customConfig.GetString("app.environment"))
	fmt.Printf("Advanced features enabled: %t\n", customConfig.GetBoolDefault("features.advanced", false))
}

// createCustomLogger creates a custom logger implementation
func createCustomLogger() logging.LoggerService {
	// For this example, we'll create a logger with NoOp implementation
	// In a real application, you might inject a structured logger like logrus, zap, etc.

	// Import the component package for ComponentConfig
	config := component.ComponentConfig{
		ID:   "custom-logger",
		Name: "Custom Logger",
		Type: component.TypeService,
	}

	// Create with NoOp logger (could be replaced with any logging.Logger implementation)
	noOpLogger := infraLogging.NewNoOpLogger()
	customLogger, err := infraLogging.NewLogger(config, noOpLogger)
	if err != nil {
		log.Fatal("Failed to create custom logger:", err)
	}

	return customLogger
}

// runCustomExamples runs all custom dependency examples
func runCustomExamples() {
	fmt.Println("=== Builder API Custom Dependencies Examples ===")
	fmt.Println("This demonstrates the key advantage of the Builder API:")
	fmt.Println("Easy injection of custom configurations, loggers, and other services")
	fmt.Println()

	// Demonstrate custom configuration
	demonstrateCustomConfig()

	// Demonstrate custom logger
	demonstrateCustomLogger()

	// Demonstrate full customization
	demonstrateFullCustomization()

	fmt.Println("\n=== All Custom Dependency Examples Completed ===")
	fmt.Println("The Builder API successfully injected custom dependencies")
	fmt.Println("without the complexity of FX dependency injection!")
}
