// Package runtime provides high-level runtime environment interfaces and implementations.
package runtime

import (
	"github.com/fintechain/skeleton/internal/domain/runtime"
	infraRuntime "github.com/fintechain/skeleton/internal/infrastructure/runtime"
)

// RuntimeEnvironment is the main interface for the runtime environment.
type RuntimeEnvironment = runtime.RuntimeEnvironment

// NewRuntime creates a new runtime environment with direct dependency injection.
var NewRuntime = infraRuntime.NewRuntime

// NewRuntimeWithOptions creates a new runtime environment with the provided options.
var NewRuntimeWithOptions = infraRuntime.NewRuntimeWithOptions

// RuntimeOption configures the runtime environment.
type RuntimeOption = infraRuntime.RuntimeOption

// WithPlugins sets the plugins to load at startup.
var WithPlugins = infraRuntime.WithPlugins

// WithRegistry sets a custom registry implementation.
var WithRegistry = infraRuntime.WithRegistry

// WithPluginManager sets a custom plugin manager implementation.
var WithPluginManager = infraRuntime.WithPluginManager

// WithEventBus sets a custom event bus implementation.
var WithEventBus = infraRuntime.WithEventBus

// WithLogger sets a custom logger implementation.
var WithLogger = infraRuntime.WithLogger

// WithConfiguration sets a custom configuration implementation.
var WithConfiguration = infraRuntime.WithConfiguration
