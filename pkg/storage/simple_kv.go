package storage

import "time"

// SimpleKV defines the kv store interface.
// SimpleKV is a simple key-value pair.

type SimpleKV interface {
	// Get gets the value for key k from the KV storage.
	// Returns a nil value if the key does not exist.
	// Key shouldn't be nil.
	Get(k []byte) ([]byte, error)
	// Set sets the value for key k as v into the KV storage.
	// Neither key nor value should be nil.
	Set(k []byte, v []byte) error
	// SetEx sets the value for key k as v into the KV storage with expiration.
	// Neither key nor value should be nil.
	SetEx(k []byte, v []byte, expiration time.Duration) error
	// HSet and HGet are based on a two-layer storage, a key k stores (key,value) pairs.
	// key k and sub-key hk and v should not be nil.
	// HSet modifies the value of hk stored under the key k, change it to v.
	HSet(k []byte, hk string, v []byte) error
	// HGet returns the value of hk under the key k.
	HGet(k []byte, hk string) ([]byte, error)
	// HGetAll returns all kv pairs under the key k
	HGetAll(k []byte) (map[string]string, error)
	// HDel deletes the value of hk under the key k.
	HDel(k []byte, hk string) error
	// HIncryBy concurrent-safely increases the value of hk by n under the key k
	HIncrBy(k []byte, hk string, n int) error
	// Exist
	Exists(k []byte) bool
	// Deletes removes the entry for key k from the KV storage.
	// Key shouldn't be nil.
	// Delete works on both basic KV storage and two-layer storage.
	Delete(k []byte) error
	// Add specific value to pointed set.
	SAdd(k []byte, v string) error
	// Remove a value from set.
	SRem(k []byte, v string) error
	// Get the number of values of set
	SCard(k []byte) (int64, error)
	// Get all values of set
	SMembers(k []byte) ([]string, error)
}
