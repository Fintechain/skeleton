package main

import (
	"log"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/system"
	pkgSystem "github.com/fintechain/skeleton/pkg/system"
)

// ExamplePlugin demonstrates a simple plugin implementation
type ExamplePlugin struct {
	id      string
	version string
}

func (p *ExamplePlugin) ID() string {
	return p.id
}

func (p *ExamplePlugin) Version() string {
	return p.version
}

func (p *ExamplePlugin) Load(ctx component.Context, registry component.Registry) error {
	log.Printf("Loading plugin: %s v%s", p.id, p.version)
	return nil
}

func (p *ExamplePlugin) Unload(ctx component.Context) error {
	log.Printf("Unloading plugin: %s v%s", p.id, p.version)
	return nil
}

func (p *ExamplePlugin) Components() []component.Component {
	return []component.Component{}
}

func main() {
	// Create configuration
	config := &system.Config{
		ServiceID: "fx-example",
		StorageConfig: storage.MultiStoreConfig{
			RootPath:      "./data",
			DefaultEngine: "memory",
		},
	}

	// Create plugins
	plugins := []plugin.Plugin{
		&ExamplePlugin{
			id:      "example-plugin",
			version: "1.0.0",
		},
	}

	// Start the system using the public API with functional options
	err := pkgSystem.StartSystem(
		pkgSystem.WithConfig(config),
		pkgSystem.WithPlugins(plugins),
	)

	if err != nil {
		log.Fatalf("Failed to start system: %v", err)
	}

	log.Println("System started successfully!")
}
