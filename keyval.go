package main

import (
	"sync"
	"time"
)

type KV struct {
	mu     sync.RWMutex         // Ensures safe concurrent access.
	data   map[string][]byte    // Map that stores key-values.
	expiry map[string]time.Time // Map that stores TTL.
}

// Initializes and returns KV
func NewKV() *KV {
	kv := &KV{
		data:   make(map[string][]byte),
		expiry: make(map[string]time.Time),
	}
	go kv.cleanupExpiredKeys() // Start the cleanup process
	return kv
}

// Uses Lock() to prevent concurrent writes.

func (kv *KV) Set(key, val []byte, ttl time.Duration) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[string(key)] = val
	if ttl > 0 {
		kv.expiry[string(key)] = time.Now().Add(ttl)
	} else {
		delete(kv.expiry, string(key)) // Remove expiry if no TTL
	}
	return nil
}

// Uses RLock() which allows multiple goroutines to read at same time.
func (kv *KV) Get(key []byte) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	if exp, ok := kv.expiry[string(key)]; ok && time.Now().After(exp) {
		delete(kv.data, string(key))
		delete(kv.expiry, string(key))
		return nil, false
	}

	val, ok := kv.data[string(key)]
	return val, ok
}

// Removal of keys from the store.
func (kv *KV) Delete(key []byte) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.data, string(key))
}

func (kv *KV) cleanupExpiredKeys() {
	for {
		time.Sleep(1 * time.Second)
		kv.mu.Lock()
		for key, exp := range kv.expiry {
			if time.Now().After(exp) {
				delete(kv.data, key)
				delete(kv.expiry, key)
			}
		}
		kv.mu.Unlock()
	}
}
