// Package memory provides an in-memory implementation of the storage engine.
package memory

import (
	"crypto/sha256"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/logging"
)

// Store implements storage.Store for in-memory storage.
type Store struct {
	name      string
	path      string
	options   Options
	data      map[string][]byte
	closed    bool
	mutex     sync.RWMutex
	versions  map[int64]map[string][]byte
	currVer   int64
	totalSize int64
	logger    logging.Logger
}

// NewStore creates a new in-memory store.
// Returns the concrete type for internal use while maintaining interface
// compliance for external API.
func NewStore(name, path string, options Options, logger logging.Logger) *Store {
	if name == "" {
		panic("store name cannot be empty")
	}

	if logger == nil {
		panic("logger cannot be nil")
	}

	logger.Debug("Creating new memory store: %s at path: %s", name, path)
	return &Store{
		name:     name,
		path:     path,
		options:  options,
		data:     make(map[string][]byte),
		versions: make(map[int64]map[string][]byte),
		currVer:  0,
		logger:   logger,
	}
}

// Get retrieves the value for the given key.
func (s *Store) Get(key []byte) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		s.logger.Debug("Attempted to access closed store: %s", s.name)
		return nil, storage.ErrStoreClosed
	}

	keyStr := string(key)
	value, exists := s.data[keyStr]
	if !exists {
		s.logger.Debug("Key not found in store %s: %s", s.name, keyStr)
		return nil, storage.ErrKeyNotFound
	}

	s.logger.Debug("Retrieved value for key in store %s: %s", s.name, keyStr)
	// Return a copy to prevent external modification
	result := make([]byte, len(value))
	copy(result, value)
	return result, nil
}

// Set stores the value for the given key.
func (s *Store) Set(key, value []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		s.logger.Debug("Attempted to set key in closed store: %s", s.name)
		return storage.ErrStoreClosed
	}

	keyStr := string(key)
	s.logger.Debug("Setting key in store %s: %s", s.name, keyStr)

	// Check if exceeding max size
	if s.options.MaxSize > 0 {
		newSize := s.totalSize
		if existing, exists := s.data[keyStr]; exists {
			newSize -= int64(len(existing))
		}
		newSize += int64(len(value))

		if newSize > s.options.MaxSize {
			s.logger.Warn("Store %s: Cannot set key %s, exceeds maximum store size of %d bytes",
				s.name, keyStr, s.options.MaxSize)
			return storage.WrapError(storage.ErrInvalidConfig, "exceeds maximum store size")
		}
		s.totalSize = newSize
	}

	// Store a copy to prevent external modification
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)
	s.data[keyStr] = valueCopy

	return nil
}

// Delete removes the key-value pair.
func (s *Store) Delete(key []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		s.logger.Debug("Attempted to delete key in closed store: %s", s.name)
		return storage.ErrStoreClosed
	}

	keyStr := string(key)
	s.logger.Debug("Deleting key in store %s: %s", s.name, keyStr)

	if existing, exists := s.data[keyStr]; exists {
		s.totalSize -= int64(len(existing))
	}
	delete(s.data, keyStr)

	return nil
}

// Has checks if a key exists.
func (s *Store) Has(key []byte) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		s.logger.Debug("Attempted to check key in closed store: %s", s.name)
		return false, storage.ErrStoreClosed
	}

	keyStr := string(key)
	_, exists := s.data[keyStr]
	s.logger.Debug("Checking key existence in store %s: %s, exists: %v", s.name, keyStr, exists)
	return exists, nil
}

// Iterate calls fn for each key-value pair.
func (s *Store) Iterate(fn func(key, value []byte) bool) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		return storage.ErrStoreClosed
	}

	for keyStr, value := range s.data {
		// Make copies to avoid modification
		key := []byte(keyStr)
		valueCopy := make([]byte, len(value))
		copy(valueCopy, value)

		if !fn(key, valueCopy) {
			break
		}
	}

	return nil
}

// Close releases resources and makes the store unusable.
func (s *Store) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		s.logger.Debug("Attempted to close already closed store: %s", s.name)
		return nil // Already closed, not an error
	}

	s.logger.Info("Closing store: %s", s.name)
	s.closed = true
	s.data = nil     // Release memory
	s.versions = nil // Release memory
	s.totalSize = 0

	return nil
}

// IsClosed returns true if the store is closed.
func (s *Store) IsClosed() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.closed
}

// Name returns the store identifier.
func (s *Store) Name() string {
	return s.name
}

// Path returns the storage path/location.
func (s *Store) Path() string {
	return s.path
}

// IterateRange iterates over keys in the specified range.
// Implements storage.RangeQueryable interface.
func (s *Store) IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		return storage.ErrStoreClosed
	}

	// Collect matching keys for sorting
	var keys []string
	startStr := string(start)
	endStr := string(end)

	for keyStr := range s.data {
		// In range check
		if (len(startStr) == 0 || keyStr >= startStr) &&
			(len(endStr) == 0 || keyStr < endStr) {
			keys = append(keys, keyStr)
		}
	}

	// Process in order
	for _, keyStr := range sortStrings(keys, ascending) {
		value := s.data[keyStr]

		// Make copies to avoid modification
		key := []byte(keyStr)
		valueCopy := make([]byte, len(value))
		copy(valueCopy, value)

		if !fn(key, valueCopy) {
			break
		}
	}

	return nil
}

// SupportsRangeQueries returns true as memory store supports range queries.
// Implements storage.RangeQueryable interface.
func (s *Store) SupportsRangeQueries() bool {
	return true
}

// SaveVersion creates a new immutable version of the store.
// Implements storage.Versioned interface.
func (s *Store) SaveVersion() (int64, []byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return 0, nil, storage.ErrStoreClosed
	}

	// Create new version number
	s.currVer++
	version := s.currVer

	// Create snapshot of current state
	snapshot := make(map[string][]byte)
	for k, v := range s.data {
		snapshot[k] = append([]byte(nil), v...)
	}

	// Store the snapshot
	s.versions[version] = snapshot

	// Clean up old versions if exceeding max versions
	if s.options.MaxVersions > 0 && len(s.versions) > s.options.MaxVersions {
		s.cleanupOldVersions()
	}

	// Calculate hash of the snapshot for integrity verification
	hash := s.calculateHash(snapshot)

	return version, hash, nil
}

// LoadVersion loads a specific version of the store.
// Implements storage.Versioned interface.
func (s *Store) LoadVersion(version int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return storage.ErrStoreClosed
	}

	snapshot, exists := s.versions[version]
	if !exists {
		return storage.ErrVersionNotFound
	}

	// Replace current data with the snapshot
	s.data = make(map[string][]byte)
	for k, v := range snapshot {
		s.data[k] = append([]byte(nil), v...)
	}

	// Recalculate size
	s.totalSize = 0
	for _, v := range s.data {
		s.totalSize += int64(len(v))
	}

	return nil
}

// ListVersions returns all available versions.
// Implements storage.Versioned interface.
func (s *Store) ListVersions() []int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	versions := make([]int64, 0, len(s.versions))
	for v := range s.versions {
		versions = append(versions, v)
	}
	return versions
}

// CurrentVersion returns the current version number.
// Implements storage.Versioned interface.
func (s *Store) CurrentVersion() int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.currVer
}

// SupportsVersioning returns true as memory store supports versioning.
// Implements storage.Versioned interface.
func (s *Store) SupportsVersioning() bool {
	return true
}

// Helper methods
func (s *Store) cleanupOldVersions() {
	// Find all versions
	versions := make([]int64, 0, len(s.versions))
	for v := range s.versions {
		versions = append(versions, v)
	}

	// Sort versions
	for i := 0; i < len(versions)-1; i++ {
		for j := i + 1; j < len(versions); j++ {
			if versions[i] > versions[j] {
				versions[i], versions[j] = versions[j], versions[i]
			}
		}
	}

	// Remove oldest versions
	toRemove := len(versions) - s.options.MaxVersions
	if toRemove <= 0 {
		return
	}

	for i := 0; i < toRemove; i++ {
		delete(s.versions, versions[i])
	}
}

func (s *Store) calculateHash(data map[string][]byte) []byte {
	h := sha256.New()

	// Sort keys for deterministic hash
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	for _, k := range sortStrings(keys, true) {
		h.Write([]byte(k))
		h.Write(data[k])
	}

	return h.Sum(nil)
}

// sortStrings returns a sorted copy of the string slice
func sortStrings(strs []string, ascending bool) []string {
	result := make([]string, len(strs))
	copy(result, strs)

	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if (ascending && result[i] > result[j]) || (!ascending && result[i] < result[j]) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}
