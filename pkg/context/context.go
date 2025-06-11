// Package context provides context interfaces for the Fintechain Skeleton framework.
package context

import (
	"github.com/fintechain/skeleton/internal/domain/context"
	infraContext "github.com/fintechain/skeleton/internal/infrastructure/context"
)

// Context is the main context interface used throughout the framework.
type Context = context.Context

// NewContext creates a new framework context instance.
var NewContext = infraContext.NewContext
