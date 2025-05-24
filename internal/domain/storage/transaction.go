// Package storage provides interfaces and types for the storage system.
package storage

// Transactional interface for stores that support transactions.
type Transactional interface {
	// BeginTx starts a new transaction.
	// Returns an error if transactions are not supported or cannot be started.
	BeginTx() (Transaction, error)

	// SupportsTransactions returns true if this store supports transactions.
	SupportsTransactions() bool
}

// Transaction represents an atomic set of operations.
// All operations within a transaction either succeed as a group or fail.
type Transaction interface {
	// Embed the Store interface - transactions support all store operations
	Store

	// Commit makes all changes permanent.
	// Returns an error if the transaction cannot be committed.
	Commit() error

	// Rollback discards all changes.
	// Returns an error if the transaction cannot be rolled back.
	Rollback() error

	// IsActive returns true if transaction is still active.
	// A transaction is active until Commit or Rollback is called.
	IsActive() bool
}
