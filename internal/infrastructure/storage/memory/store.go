// Package memory provides an in-memory storage store implementation.
package memory

import (
	"errors"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/storage"
)

// Store implements the storage.Store interface for in-memory storage.
type Store struct {
	name string
	path string
	data map[string][]byte
	mu   sync.RWMutex
}

// NewStore creates a new in-memory store instance.
func NewStore(name, path string) *Store {
	return &Store{
		name: name,
		path: path,
		data: make(map[string][]byte),
	}
}

// Get retrieves the value associated with the given key.
func (s *Store) Get(key []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[string(key)]
	if !exists {
		return nil, errors.New(storage.ErrKeyNotFound)
	}

	// Return a copy to prevent external modification
	result := make([]byte, len(value))
	copy(result, value)
	return result, nil
}

// Set stores a value for the given key.
func (s *Store) Set(key, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store a copy to prevent external modification
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)
	s.data[string(key)] = valueCopy

	return nil
}

// Delete removes the key-value pair for the given key.
func (s *Store) Delete(key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyStr := string(key)
	if _, exists := s.data[keyStr]; !exists {
		return errors.New(storage.ErrKeyNotFound)
	}

	delete(s.data, keyStr)
	return nil
}

// Has checks whether a key exists in the store.
func (s *Store) Has(key []byte) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.data[string(key)]
	return exists, nil
}

// Iterate calls the provided function for each key-value pair in the store.
func (s *Store) Iterate(fn func(key, value []byte) bool) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for k, v := range s.data {
		// Create copies to prevent external modification
		keyCopy := []byte(k)
		valueCopy := make([]byte, len(v))
		copy(valueCopy, v)

		if !fn(keyCopy, valueCopy) {
			break
		}
	}

	return nil
}

// Close releases all resources associated with the store.
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = nil
	return nil
}

// Name returns the name of this store instance.
func (s *Store) Name() string {
	return s.name
}

// Path returns the storage path or location for this store.
func (s *Store) Path() string {
	return s.path
}
