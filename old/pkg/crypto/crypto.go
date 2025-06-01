// Package crypto provides cryptographic utilities for the skeleton framework.
package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// ===== HASH FUNCTIONS =====

// HashSHA256 calculates the SHA-256 hash of the input string and returns the hexadecimal representation.
//
// Example usage:
//
//	hash := crypto.HashSHA256("hello world")
//	fmt.Println(hash) // outputs: b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
func HashSHA256(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// HashSHA256Bytes calculates the SHA-256 hash of the input bytes and returns the hexadecimal representation.
//
// Example usage:
//
//	hash := crypto.HashSHA256Bytes([]byte("hello world"))
func HashSHA256Bytes(input []byte) string {
	hasher := sha256.New()
	hasher.Write(input)
	return hex.EncodeToString(hasher.Sum(nil))
}

// HashSHA1 calculates the SHA-1 hash of the input string and returns the hexadecimal representation.
//
// Example usage:
//
//	hash := crypto.HashSHA1("hello world")
func HashSHA1(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// HashSHA512 calculates the SHA-512 hash of the input string and returns the hexadecimal representation.
//
// Example usage:
//
//	hash := crypto.HashSHA512("hello world")
func HashSHA512(input string) string {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// HashMD5 calculates the MD5 hash of the input string and returns the hexadecimal representation.
// Note: MD5 is cryptographically broken and should not be used for security purposes.
//
// Example usage:
//
//	hash := crypto.HashMD5("hello world")
func HashMD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ===== HMAC FUNCTIONS =====

// HMACSHA256 calculates the HMAC-SHA256 of the input string using the provided key.
//
// Example usage:
//
//	hmac := crypto.HMACSHA256("secret-key", "hello world")
func HMACSHA256(key, input string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA256Bytes calculates the HMAC-SHA256 of the input bytes using the provided key.
//
// Example usage:
//
//	hmac := crypto.HMACSHA256Bytes([]byte("secret-key"), []byte("hello world"))
func HMACSHA256Bytes(key, input []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(input)
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA512 calculates the HMAC-SHA512 of the input string using the provided key.
//
// Example usage:
//
//	hmac := crypto.HMACSHA512("secret-key", "hello world")
func HMACSHA512(key, input string) string {
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// ===== RANDOM GENERATION =====

// GenerateRandomBytes generates cryptographically secure random bytes of the specified length.
//
// Example usage:
//
//	randomBytes, err := crypto.GenerateRandomBytes(32)
//	if err != nil {
//	    // Handle error
//	}
func GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, NewError(ErrInvalidInput, "length must be positive", nil)
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, NewError(ErrRandomGeneration, "failed to generate random bytes", err)
	}
	return bytes, nil
}

// GenerateRandomString generates a cryptographically secure random string of the specified length.
// The string contains only alphanumeric characters (a-z, A-Z, 0-9).
//
// Example usage:
//
//	randomString, err := crypto.GenerateRandomString(16)
//	if err != nil {
//	    // Handle error
//	}
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	if length <= 0 {
		return "", NewError(ErrInvalidInput, "length must be positive", nil)
	}

	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	result := make([]byte, length)
	for i, b := range bytes {
		result[i] = charset[int(b)%len(charset)]
	}

	return string(result), nil
}

// GenerateRandomHex generates a cryptographically secure random hexadecimal string.
//
// Example usage:
//
//	randomHex, err := crypto.GenerateRandomHex(16) // generates 32-character hex string
//	if err != nil {
//	    // Handle error
//	}
func GenerateRandomHex(byteLength int) (string, error) {
	bytes, err := GenerateRandomBytes(byteLength)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ===== ENCODING UTILITIES =====

// EncodeBase64 encodes the input string to base64.
//
// Example usage:
//
//	encoded := crypto.EncodeBase64("hello world")
func EncodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// DecodeBase64 decodes a base64 string.
//
// Example usage:
//
//	decoded, err := crypto.DecodeBase64("aGVsbG8gd29ybGQ=")
//	if err != nil {
//	    // Handle error
//	}
func DecodeBase64(input string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", NewError(ErrInvalidInput, "invalid base64 input", err)
	}
	return string(decoded), nil
}

// EncodeHex encodes the input string to hexadecimal.
//
// Example usage:
//
//	encoded := crypto.EncodeHex("hello world")
func EncodeHex(input string) string {
	return hex.EncodeToString([]byte(input))
}

// DecodeHex decodes a hexadecimal string.
//
// Example usage:
//
//	decoded, err := crypto.DecodeHex("68656c6c6f20776f726c64")
//	if err != nil {
//	    // Handle error
//	}
func DecodeHex(input string) (string, error) {
	decoded, err := hex.DecodeString(input)
	if err != nil {
		return "", NewError(ErrInvalidInput, "invalid hex input", err)
	}
	return string(decoded), nil
}

// ===== VALIDATION UTILITIES =====

// VerifyHMAC verifies an HMAC signature against the expected value.
//
// Example usage:
//
//	isValid := crypto.VerifyHMAC("secret-key", "hello world", expectedHMAC)
func VerifyHMAC(key, input, expectedHMAC string) bool {
	computedHMAC := HMACSHA256(key, input)
	return hmac.Equal([]byte(computedHMAC), []byte(expectedHMAC))
}

// VerifyHMACBytes verifies an HMAC signature against the expected value using byte inputs.
//
// Example usage:
//
//	isValid := crypto.VerifyHMACBytes(keyBytes, inputBytes, expectedHMACBytes)
func VerifyHMACBytes(key, input, expectedHMAC []byte) bool {
	h := hmac.New(sha256.New, key)
	h.Write(input)
	computedHMAC := h.Sum(nil)
	return hmac.Equal(computedHMAC, expectedHMAC)
}

// ===== HASH UTILITIES =====

// HashWithAlgorithm calculates a hash using the specified algorithm.
//
// Example usage:
//
//	hash, err := crypto.HashWithAlgorithm("sha256", "hello world")
//	if err != nil {
//	    // Handle error
//	}
func HashWithAlgorithm(algorithm, input string) (string, error) {
	var hasher hash.Hash

	switch algorithm {
	case "md5":
		hasher = md5.New()
	case "sha1":
		hasher = sha1.New()
	case "sha256":
		hasher = sha256.New()
	case "sha512":
		hasher = sha512.New()
	default:
		return "", NewError(ErrUnsupportedAlgorithm, fmt.Sprintf("unsupported hash algorithm: %s", algorithm), nil)
	}

	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// ===== ERROR CONSTANTS =====

// Common crypto error codes
const (
	ErrInvalidInput         = "crypto.invalid_input"
	ErrRandomGeneration     = "crypto.random_generation_failed"
	ErrUnsupportedAlgorithm = "crypto.unsupported_algorithm"
	ErrEncodingFailed       = "crypto.encoding_failed"
	ErrDecodingFailed       = "crypto.decoding_failed"
)

// ===== ERROR HANDLING =====

// Error represents a domain-specific error from the crypto system.
type Error = component.Error

// NewError creates a new crypto error with the given code, message, and optional cause.
func NewError(code, message string, cause error) *Error {
	return component.NewError(code, message, cause)
}

// IsCryptoError checks if an error is a crypto error with the given code.
func IsCryptoError(err error, code string) bool {
	return component.IsComponentError(err, code)
}
