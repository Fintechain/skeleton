package registry

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fintechain/skeleton/pkg/registry"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// TestNewRegistry tests the registry constructor function
func TestNewRegistry(t *testing.T) {
	reg := registry.NewRegistry()

	assert.NotNil(t, reg)
	assert.Equal(t, 0, reg.Count())
	assert.Empty(t, reg.List())
}

// TestRegistryInterfaceCompliance verifies the implementation satisfies the domain interface
func TestRegistryInterfaceCompliance(t *testing.T) {
	// Verify interface compliance
	var _ registry.Registry = registry.NewRegistry()
}

// TestRegistryOperations tests all registry CRUD operations
func TestRegistryOperations(t *testing.T) {
	tests := []struct {
		name        string
		operation   func(registry.Registry, *mocks.Factory) error
		expectError bool
		description string
	}{
		{
			name: "register valid item",
			operation: func(r registry.Registry, f *mocks.Factory) error {
				item := f.ComponentInterface().(*mocks.MockComponent)
				item.SetID("test-item")
				return r.Register(item)
			},
			expectError: false,
			description: "Should register valid items successfully",
		},
		{
			name: "register nil item",
			operation: func(r registry.Registry, f *mocks.Factory) error {
				return r.Register(nil)
			},
			expectError: true,
			description: "Should reject nil items",
		},
		{
			name: "register item with empty ID",
			operation: func(r registry.Registry, f *mocks.Factory) error {
				item := f.ComponentInterface().(*mocks.MockComponent)
				item.SetID("")
				return r.Register(item)
			},
			expectError: true,
			description: "Should reject items with empty IDs",
		},
		{
			name: "register duplicate item",
			operation: func(r registry.Registry, f *mocks.Factory) error {
				item := f.ComponentInterface().(*mocks.MockComponent)
				item.SetID("duplicate-item")
				r.Register(item)
				return r.Register(item)
			},
			expectError: true,
			description: "Should reject duplicate items",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := registry.NewRegistry()
			factory := mocks.NewFactory()

			err := tt.operation(reg, factory)

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

// TestRegistryGet tests item retrieval functionality
func TestRegistryGet(t *testing.T) {
	reg := registry.NewRegistry()
	factory := mocks.NewFactory()

	// Test getting non-existent item
	_, err := reg.Get("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), registry.ErrItemNotFound)

	// Register an item
	item := factory.ComponentInterface().(*mocks.MockComponent)
	item.SetID("test-item")
	item.SetName("Test Item")
	err = reg.Register(item)
	require.NoError(t, err)

	// Test getting existing item
	retrieved, err := reg.Get("test-item")
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "test-item", retrieved.ID())
	assert.Equal(t, "Test Item", retrieved.Name())
}

// TestRegistryHas tests item existence checking
func TestRegistryHas(t *testing.T) {
	reg := registry.NewRegistry()
	factory := mocks.NewFactory()

	// Test non-existent item
	assert.False(t, reg.Has("non-existent"))

	// Register an item
	item := factory.ComponentInterface().(*mocks.MockComponent)
	item.SetID("test-item")
	err := reg.Register(item)
	require.NoError(t, err)

	// Test existing item
	assert.True(t, reg.Has("test-item"))
}

// TestRegistryRemove tests item removal functionality
func TestRegistryRemove(t *testing.T) {
	reg := registry.NewRegistry()
	factory := mocks.NewFactory()

	// Test removing non-existent item
	err := reg.Remove("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), registry.ErrItemNotFound)

	// Register an item
	item := factory.ComponentInterface().(*mocks.MockComponent)
	item.SetID("test-item")
	err = reg.Register(item)
	require.NoError(t, err)

	// Verify item exists
	assert.True(t, reg.Has("test-item"))
	assert.Equal(t, 1, reg.Count())

	// Remove the item
	err = reg.Remove("test-item")
	assert.NoError(t, err)

	// Verify item is removed
	assert.False(t, reg.Has("test-item"))
	assert.Equal(t, 0, reg.Count())
}

// TestRegistryList tests listing all items
func TestRegistryList(t *testing.T) {
	reg := registry.NewRegistry()
	factory := mocks.NewFactory()

	// Test empty registry
	items := reg.List()
	assert.Empty(t, items)

	// Register multiple items
	item1 := factory.ComponentInterface().(*mocks.MockComponent)
	item1.SetID("item-1")
	item1.SetName("Item 1")

	item2 := factory.ComponentInterface().(*mocks.MockComponent)
	item2.SetID("item-2")
	item2.SetName("Item 2")

	err := reg.Register(item1)
	require.NoError(t, err)
	err = reg.Register(item2)
	require.NoError(t, err)

	// Test listing items
	items = reg.List()
	assert.Len(t, items, 2)

	// Verify items are in the list (order may vary)
	ids := make(map[string]bool)
	for _, item := range items {
		ids[item.ID()] = true
	}
	assert.True(t, ids["item-1"])
	assert.True(t, ids["item-2"])
}

// TestRegistryCount tests item counting
func TestRegistryCount(t *testing.T) {
	reg := registry.NewRegistry()
	factory := mocks.NewFactory()

	// Test empty registry
	assert.Equal(t, 0, reg.Count())

	// Register items and verify count
	for i := 0; i < 5; i++ {
		item := factory.ComponentInterface().(*mocks.MockComponent)
		item.SetID(fmt.Sprintf("item-%d", i))
		err := reg.Register(item)
		require.NoError(t, err)
		assert.Equal(t, i+1, reg.Count())
	}

	// Remove items and verify count
	for i := 0; i < 3; i++ {
		err := reg.Remove(fmt.Sprintf("item-%d", i))
		require.NoError(t, err)
		assert.Equal(t, 5-(i+1), reg.Count())
	}
}

// TestRegistryClear tests clearing all items
func TestRegistryClear(t *testing.T) {
	reg := registry.NewRegistry()
	factory := mocks.NewFactory()

	// Register multiple items
	for i := 0; i < 5; i++ {
		item := factory.ComponentInterface().(*mocks.MockComponent)
		item.SetID(fmt.Sprintf("item-%d", i))
		err := reg.Register(item)
		require.NoError(t, err)
	}

	// Verify items are registered
	assert.Equal(t, 5, reg.Count())

	// Clear registry
	reg.Clear()

	// Verify registry is empty
	assert.Equal(t, 0, reg.Count())
	assert.Empty(t, reg.List())
	assert.False(t, reg.Has("item-0"))
}

// TestRegistryErrorHandling tests error conditions with string constants
func TestRegistryErrorHandling(t *testing.T) {
	reg := registry.NewRegistry()
	factory := mocks.NewFactory()

	// Test item not found error
	_, err := reg.Get("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), registry.ErrItemNotFound)

	// Test remove non-existent item error
	err = reg.Remove("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), registry.ErrItemNotFound)

	// Test invalid item error (nil)
	err = reg.Register(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), registry.ErrInvalidItem)

	// Test invalid item error (empty ID)
	item := factory.ComponentInterface().(*mocks.MockComponent)
	item.SetID("")
	err = reg.Register(item)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), registry.ErrInvalidItem)

	// Test duplicate item error
	validItem := factory.ComponentInterface().(*mocks.MockComponent)
	validItem.SetID("test-item")
	err = reg.Register(validItem)
	require.NoError(t, err)

	err = reg.Register(validItem)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), registry.ErrItemAlreadyExists)
}

// TestRegistryConcurrency tests thread-safety with concurrent operations
func TestRegistryConcurrency(t *testing.T) {
	reg := registry.NewRegistry()
	factory := mocks.NewFactory()

	var wg sync.WaitGroup
	numGoroutines := 10
	itemsPerGoroutine := 10

	// Concurrent registrations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				item := factory.ComponentInterface().(*mocks.MockComponent)
				item.SetID(fmt.Sprintf("item-%d-%d", goroutineID, j))
				item.SetName(fmt.Sprintf("Item %d-%d", goroutineID, j))
				err := reg.Register(item)
				assert.NoError(t, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify all items were registered
	expectedCount := numGoroutines * itemsPerGoroutine
	assert.Equal(t, expectedCount, reg.Count())

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				id := fmt.Sprintf("item-%d-%d", goroutineID, j)
				item, err := reg.Get(id)
				assert.NoError(t, err)
				assert.NotNil(t, item)
				assert.Equal(t, id, item.ID())
			}
		}(i)
	}

	wg.Wait()

	// Concurrent removals
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine/2; j++ {
				id := fmt.Sprintf("item-%d-%d", goroutineID, j)
				err := reg.Remove(id)
				assert.NoError(t, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify correct number of items remain
	expectedRemaining := numGoroutines * (itemsPerGoroutine - itemsPerGoroutine/2)
	assert.Equal(t, expectedRemaining, reg.Count())
}

// TestRegistryMixedOperations tests various operations in combination
func TestRegistryMixedOperations(t *testing.T) {
	reg := registry.NewRegistry()
	factory := mocks.NewFactory()

	// Register initial items
	for i := 0; i < 5; i++ {
		item := factory.ComponentInterface().(*mocks.MockComponent)
		item.SetID(fmt.Sprintf("item-%d", i))
		item.SetName(fmt.Sprintf("Item %d", i))
		item.SetDescription(fmt.Sprintf("Description for item %d", i))
		item.SetVersion("1.0.0")
		err := reg.Register(item)
		require.NoError(t, err)
	}

	// Test mixed operations
	assert.Equal(t, 5, reg.Count())
	assert.True(t, reg.Has("item-2"))

	// Remove some items
	err := reg.Remove("item-1")
	assert.NoError(t, err)
	err = reg.Remove("item-3")
	assert.NoError(t, err)

	assert.Equal(t, 3, reg.Count())
	assert.False(t, reg.Has("item-1"))
	assert.False(t, reg.Has("item-3"))

	// Add new items
	newItem := factory.ComponentInterface().(*mocks.MockComponent)
	newItem.SetID("new-item")
	newItem.SetName("New Item")
	err = reg.Register(newItem)
	assert.NoError(t, err)

	assert.Equal(t, 4, reg.Count())
	assert.True(t, reg.Has("new-item"))

	// Verify remaining items
	items := reg.List()
	assert.Len(t, items, 4)

	expectedIDs := map[string]bool{
		"item-0":   true,
		"item-2":   true,
		"item-4":   true,
		"new-item": true,
	}

	for _, item := range items {
		assert.True(t, expectedIDs[item.ID()], "Unexpected item ID: %s", item.ID())
	}
}
