/*
Package storage implements the persistent layer of deepshare backend.
*/
package storage

// KV defines the kv store interface.
// kv is a simple key-value pair.
// TODO: what guarantees KV supports?
type KV interface {
	// Begin starts a new transaction.
	Begin(writable bool) (Tx, error)
	// Close the KV
	Close() error
}

// Tx defines the interface for operations inside a Transaction.
// A Tx can only be used inside one go routine.
type Tx interface {
	// Get gets the value for key k from the KV storage.
	// Returns a nil value if the key does not exist.
	// Key shouldn't be nil.
	Get(k []byte) []byte
	// Set sets the value for key k as v into the KV storage.
	// Neither key nor value should be nil.
	Set(k []byte, v []byte)
	// Deletes removes the entry for key k from the KV storage.
	// Key shouldn't be nil.
	Delete(k []byte)
	// Commit commits the transaction operations to the KV storage.
	Commit() error
	// Rollback undoes the transaction operations to the KV storage.
	Rollback() error
}
