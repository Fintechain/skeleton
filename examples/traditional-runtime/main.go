// Package main demonstrates traditional runtime usage patterns.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/plugin"
	"github.com/fintechain/skeleton/pkg/runtime"
)

// Example Operation: Simple Calculator
type SimpleCalculator struct {
	*component.BaseOperation
}

func NewSimpleCalculator() *SimpleCalculator {
	config := component.ComponentConfig{
		ID:          "simple-calculator",
		Name:        "Simple Calculator",
		Description: "Performs basic arithmetic operations using traditional runtime",
		Version:     "1.0.0",
	}
	return &SimpleCalculator{
		BaseOperation: component.NewBaseOperation(config),
	}
}

func (c *SimpleCalculator) Execute(ctx context.Context, input component.Input) (component.Output, error) {
	data, ok := input.Data.(map[string]interface{})
	if !ok {
		return component.Output{}, fmt.Errorf("invalid input data")
	}

	operation, _ := data["operation"].(string)
	a, _ := data["a"].(float64)
	b, _ := data["b"].(float64)

	var result float64
	switch operation {
	case "add":
		result = a + b
	case "multiply":
		result = a * b
	default:
		return component.Output{}, fmt.Errorf("unknown operation: %s", operation)
	}

	return component.Output{
		Data: map[string]interface{}{
			"result": result,
		},
	}, nil
}

// Example Plugin using traditional patterns
type SimpleCalculatorPlugin struct {
	*component.BaseService
	calculator *SimpleCalculator
}

func NewSimpleCalculatorPlugin() *SimpleCalculatorPlugin {
	config := component.ComponentConfig{
		ID:          "simple-calculator-plugin",
		Name:        "Simple Calculator Plugin",
		Description: "Provides basic calculator operations",
		Version:     "1.0.0",
	}
	return &SimpleCalculatorPlugin{
		BaseService: component.NewBaseService(config),
		calculator:  NewSimpleCalculator(),
	}
}

func (p *SimpleCalculatorPlugin) Author() string {
	return "Fintechain Team"
}

func (p *SimpleCalculatorPlugin) PluginType() plugin.PluginType {
	return plugin.TypeProcessor
}

func (p *SimpleCalculatorPlugin) Initialize(ctx context.Context, system component.System) error {
	if err := p.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Register the calculator operation
	registry := system.Registry()
	if err := registry.Register(p.calculator); err != nil {
		return fmt.Errorf("failed to register calculator: %w", err)
	}

	fmt.Println("[SimpleCalculatorPlugin] Calculator registered")
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Traditional Runtime Demo - Fintechain Skeleton Framework")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  go run examples/traditional-runtime/main.go daemon     # Run as daemon")
		fmt.Println("  go run examples/traditional-runtime/main.go command    # Run command mode")
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "daemon":
		runTraditionalDaemon()
	case "command":
		runTraditionalCommand()
	default:
		fmt.Printf("Unknown mode: %s\n", mode)
		os.Exit(1)
	}
}

func runTraditionalDaemon() {
	fmt.Println("=== Traditional Runtime - Daemon Mode ===")

	// Create runtime using traditional options pattern
	rt, err := runtime.NewRuntimeWithOptions(
		runtime.WithPlugins(NewSimpleCalculatorPlugin()),
	)
	if err != nil {
		log.Fatalf("Failed to create runtime: %v", err)
	}

	// Start the runtime
	ctx := context.NewContext()
	if err := rt.Start(ctx); err != nil {
		log.Fatalf("Failed to start runtime: %v", err)
	}

	fmt.Println("Runtime started successfully. Press Ctrl+C to exit.")

	// In a real daemon, you'd block here until shutdown signal
	// For demo purposes, we'll just show it started
	fmt.Println("Daemon running... (in real app, this would block until shutdown)")

	// Clean shutdown
	if err := rt.Stop(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}

func runTraditionalCommand() {
	fmt.Println("=== Traditional Runtime - Command Mode ===")

	// Create runtime
	rt, err := runtime.NewRuntimeWithOptions(
		runtime.WithPlugins(NewSimpleCalculatorPlugin()),
	)
	if err != nil {
		log.Fatalf("Failed to create runtime: %v", err)
	}

	// Start runtime
	ctx := context.NewContext()
	if err := rt.Start(ctx); err != nil {
		log.Fatalf("Failed to start runtime: %v", err)
	}
	defer rt.Stop(ctx)

	// Execute operations using traditional runtime
	fmt.Println("\n--- Addition: 15 + 25 ---")
	result, err := rt.ExecuteOperation(ctx, "simple-calculator", component.Input{
		Data: map[string]interface{}{
			"operation": "add",
			"a":         15.0,
			"b":         25.0,
		},
		Metadata: map[string]string{"mode": "traditional"},
	})

	if err != nil {
		fmt.Printf("Addition failed: %v\n", err)
	} else {
		fmt.Printf("Result: %v\n", result.Data)
	}

	// Execute multiplication
	fmt.Println("\n--- Multiplication: 6 * 7 ---")
	result, err = rt.ExecuteOperation(ctx, "simple-calculator", component.Input{
		Data: map[string]interface{}{
			"operation": "multiply",
			"a":         6.0,
			"b":         7.0,
		},
		Metadata: map[string]string{"mode": "traditional"},
	})

	if err != nil {
		fmt.Printf("Multiplication failed: %v\n", err)
	} else {
		fmt.Printf("Result: %v\n", result.Data)
	}

	fmt.Println("\n=== Traditional Runtime Examples Complete ===")
}
