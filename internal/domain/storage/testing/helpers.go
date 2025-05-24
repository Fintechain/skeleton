// Package testing provides compliance tests for storage implementations.
package testing

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/stretchr/testify/require"
)

// randomSource is the shared random source for test data generation
var randomSource = rand.NewSource(time.Now().UnixNano())
var random = rand.New(randomSource)

// GenerateTestData generates random key-value pairs for testing.
// It returns two slices: keys and values.
func GenerateTestData(numPairs int) ([][]byte, [][]byte) {
	keys := make([][]byte, numPairs)
	values := make([][]byte, numPairs)

	for i := 0; i < numPairs; i++ {
		keySize := 8 + random.Intn(16)    // 8-24 byte keys
		valueSize := 16 + random.Intn(64) // 16-80 byte values

		keys[i] = make([]byte, keySize)
		values[i] = make([]byte, valueSize)

		random.Read(keys[i])
		random.Read(values[i])
	}

	return keys, values
}

// PopulateStore inserts test data into a store.
func PopulateStore(t *testing.T, store storage.Store, numPairs int) ([][]byte, [][]byte) {
	keys, values := GenerateTestData(numPairs)

	for i := 0; i < numPairs; i++ {
		err := store.Set(keys[i], values[i])
		require.NoError(t, err)
	}

	return keys, values
}

// VerifyStoreContents checks that a store contains the expected key-value pairs.
func VerifyStoreContents(t *testing.T, store storage.Store, keys, values [][]byte) {
	require.Equal(t, len(keys), len(values), "keys and values must have the same length")

	for i := 0; i < len(keys); i++ {
		value, err := store.Get(keys[i])
		require.NoError(t, err)
		require.Equal(t, values[i], value)
	}
}

// BenchmarkStoreOperations runs performance benchmarks on a store.
func BenchmarkStoreOperations(b *testing.B, store storage.Store) {
	// Generate test data
	keys, values := GenerateTestData(1000)

	// Benchmark Set
	b.Run("Set", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			idx := i % len(keys)
			_ = store.Set(keys[idx], values[idx])
		}
	})

	// Populate store for Get benchmark
	for i := 0; i < len(keys); i++ {
		_ = store.Set(keys[i], values[i])
	}

	// Benchmark Get
	b.Run("Get", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			idx := i % len(keys)
			_, _ = store.Get(keys[idx])
		}
	})

	// Benchmark Has
	b.Run("Has", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			idx := i % len(keys)
			_, _ = store.Has(keys[idx])
		}
	})

	// Benchmark Delete
	b.Run("Delete", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			idx := i % len(keys)
			_ = store.Delete(keys[idx])
			_ = store.Set(keys[idx], values[idx]) // Restore for next iteration
		}
	})
}

// TestRangeQueryable tests range query functionality.
func TestRangeQueryable(t *testing.T, store storage.Store) {
	rq, ok := store.(storage.RangeQueryable)
	if !ok {
		t.Skip("Store does not support range queries")
	}

	// Populate with sequential keys for range testing
	keyPrefix := []byte("range-test-")
	for i := 0; i < 100; i++ {
		key := append([]byte(nil), keyPrefix...)
		key = append(key, byte(i))
		value := []byte{byte(i)}

		err := store.Set(key, value)
		require.NoError(t, err)
	}

	// Test ascending range query
	start := append([]byte(nil), keyPrefix...)
	start = append(start, byte(20))

	end := append([]byte(nil), keyPrefix...)
	end = append(end, byte(30))

	count := 0
	err := rq.IterateRange(start, end, true, func(key, value []byte) bool {
		require.GreaterOrEqual(t, key[len(keyPrefix)], byte(20))
		require.Less(t, key[len(keyPrefix)], byte(30))
		require.Equal(t, value[0], key[len(keyPrefix)])
		count++
		return true
	})

	require.NoError(t, err)
	require.Equal(t, 10, count)
}
