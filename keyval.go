package main

import "sync"

type KV struct {
	mu   sync.RWMutex      // Ensures safe concurrent access.
	data map[string][]byte // Map that stores key-values.
}

// Initializes and returns KV
func NewKV() *KV {
	return &KV{
		data: map[string][]byte{},
	}
}

// Uses Lock() to prevent concurrent writes.

func (kv *KV) Set(key, val []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[string(key)] = val
	return nil
}

// Uses RLock() which allows multiple goroutines to read at same time.
func (kv *KV) Get(key []byte) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	val, ok := kv.data[string(key)]
	return val, ok
}

// Removal of keys from the store.
func (kv *KV) Delete(key []byte) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.data, string(key))
}
