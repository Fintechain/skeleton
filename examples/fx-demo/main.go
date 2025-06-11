// Package main demonstrates FX integration usage for both daemon and command modes.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/fx"
	"github.com/fintechain/skeleton/pkg/plugin"
)

// Example Operation: Calculator
type Calculator struct {
	*component.BaseOperation
}

func NewCalculator() *Calculator {
	config := component.ComponentConfig{
		ID:          "calculator",
		Name:        "Calculator",
		Description: "Performs basic arithmetic operations",
		Version:     "1.0.0",
	}
	return &Calculator{
		BaseOperation: component.NewBaseOperation(config),
	}
}

func (c *Calculator) Execute(ctx context.Context, input component.Input) (component.Output, error) {
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
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return component.Output{}, fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return component.Output{}, fmt.Errorf("unknown operation: %s", operation)
	}

	return component.Output{
		Data: map[string]interface{}{
			"result":    result,
			"operation": operation,
			"operands":  []float64{a, b},
		},
	}, nil
}

// Example Plugin: Calculator Plugin
type CalculatorPlugin struct {
	*component.BaseService
	calculator *Calculator
}

func NewCalculatorPlugin() *CalculatorPlugin {
	config := component.ComponentConfig{
		ID:          "calculator-plugin",
		Name:        "Calculator Plugin",
		Description: "Provides calculator operations",
		Version:     "1.0.0",
	}
	return &CalculatorPlugin{
		BaseService: component.NewBaseService(config),
		calculator:  NewCalculator(),
	}
}

func (p *CalculatorPlugin) Author() string {
	return "Fintechain Team"
}

func (p *CalculatorPlugin) PluginType() plugin.PluginType {
	return plugin.TypeProcessor
}

func (p *CalculatorPlugin) Initialize(ctx context.Context, system component.System) error {
	if err := p.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Register the calculator operation with the system
	registry := system.Registry()
	if err := registry.Register(p.calculator); err != nil {
		return fmt.Errorf("failed to register calculator: %w", err)
	}

	fmt.Println("[CalculatorPlugin] Calculator operation registered")
	return nil
}

func (p *CalculatorPlugin) Start(ctx context.Context) error {
	fmt.Println("[CalculatorPlugin] Calculator plugin started")
	return p.BaseService.Start(ctx)
}

func (p *CalculatorPlugin) Stop(ctx context.Context) error {
	fmt.Println("[CalculatorPlugin] Calculator plugin stopped")
	return p.BaseService.Stop(ctx)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("FX Integration Demo - Fintechain Skeleton Framework")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  go run examples/fx-demo/main.go daemon     # Run as daemon")
		fmt.Println("  go run examples/fx-demo/main.go command    # Run command mode")
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

func runDaemonMode() {
	fmt.Println("=== FX Daemon Mode Example ===")

	// Create plugins
	calculatorPlugin := NewCalculatorPlugin()

	// Start daemon with plugins
	err := fx.StartDaemon(
		fx.WithPlugins(calculatorPlugin),
	)

	if err != nil {
		log.Fatalf("Daemon failed: %v", err)
	}
}

func runCommandMode() {
	fmt.Println("=== FX Command Mode Example ===")

	// Create calculator plugin
	calculatorPlugin := NewCalculatorPlugin()

	// Test addition
	fmt.Println("\n--- Addition: 10 + 5 ---")
	result, err := fx.ExecuteCommand("calculator", map[string]interface{}{
		"operation": "add",
		"a":         10.0,
		"b":         5.0,
	}, fx.WithPlugins(calculatorPlugin))

	if err != nil {
		fmt.Printf("Addition failed: %v\n", err)
	} else {
		fmt.Printf("Result: %v\n", result)
	}

	// Test division
	fmt.Println("\n--- Division: 20 / 4 ---")
	result, err = fx.ExecuteCommand("calculator", map[string]interface{}{
		"operation": "divide",
		"a":         20.0,
		"b":         4.0,
	}, fx.WithPlugins(calculatorPlugin))

	if err != nil {
		fmt.Printf("Division failed: %v\n", err)
	} else {
		fmt.Printf("Result: %v\n", result)
	}

	fmt.Println("\n=== Command Mode Examples Complete ===")
}
