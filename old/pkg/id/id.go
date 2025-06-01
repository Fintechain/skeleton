// Package id provides ID generation utilities for the skeleton framework.
package id

import (
	"encoding/binary"
	"fmt"

	"github.com/fintechain/skeleton/pkg/crypto"
)

// IDGeneratorInterface interface defines the behavior of a process ID generator.
type IDGeneratorInterface interface {
	GenerateID() (string, error)
}

// ProcessIDGenerator provides functionality to generate unique process IDs.
type ProcessIDGenerator struct {
	IDGeneratorInterface
	prefix string
}

// NewProcessIDGenerator creates a new instance of ProcessIDGenerator with the given prefix.
func NewProcessIDGenerator(prefix string) *ProcessIDGenerator {
	return &ProcessIDGenerator{
		prefix: prefix,
	}
}

// GenerateID generates a unique process ID.
func (gen *ProcessIDGenerator) GenerateID() (string, error) {
	// Check if the prefix is empty
	if gen.prefix == "" {
		return "", fmt.Errorf("prefix cannot be empty")
	}

	// Generate random bytes using the crypto package
	randomBytes, err := crypto.GenerateRandomBytes(4) // 4 bytes = 32-bit number
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert bytes to a number (0 to 4294967295, but we'll mod it to keep it reasonable)
	randomNum := binary.BigEndian.Uint32(randomBytes) % 1000000

	// Combine prefix and random number to create the process ID
	processID := fmt.Sprintf("%s-%d", gen.prefix, randomNum)

	return processID, nil
}
