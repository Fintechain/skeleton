// skeleton-old/pkg/common/id.go
package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// IDGeneratorInterface interface defines the behavior of a process ID generator.
type IDGeneratorInterface interface {
	GenerateID() (string, error)
}

// RandomProcessIDGenerator provides functionality to generate unique process IDs.
type ProcessIDGenerator struct {
	IDGeneratorInterface
	prefix string
}

// NewRandomProcessIDGenerator creates a new instance of RandomProcessIDGenerator with the given prefix.
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

	// Generate a random number to ensure uniqueness
	randomNum, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	// Combine prefix and random number to create the process ID
	processID := fmt.Sprintf("%s-%d", gen.prefix, randomNum)

	return processID, nil
}

// hashSHA256 calculates the SHA-256 hash of the input string and returns the hexadecimal representation.
func HashSHA256(input string) string {
	// Create a new SHA-256 hasher
	hasher := sha256.New()

	// Write the input string to the hasher
	hasher.Write([]byte(input))

	// Calculate the hash sum
	hashSum := hasher.Sum(nil)

	// Convert the hash sum to a hexadecimal string
	hashHex := hex.EncodeToString(hashSum)

	return hashHex
}
