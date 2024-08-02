package db

import (
	"sync"
)

// HashMap represents a simple in-memory hash table
type HashMap struct {
	mu   sync.RWMutex
	data map[string]map[string]string
}

// NewHashMap creates a new HashTable
func NewHashMap() *HashMap {
	return &HashMap{
		data: make(map[string]map[string]string),
	}
}

// HSet sets the value for a field in a hash
func (h *HashMap) HSet(key string, fields ...string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.data[key]; !ok {
		h.data[key] = make(map[string]string)
	}

	for i := 0; i < len(fields); i += 2 {
		h.data[key][fields[i]] = fields[i+1]
	}
	return true
}

// HGet gets the value of a field in a hash
func (h *HashMap) HGet(key, field string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if fields, ok := h.data[key]; ok {
		if value, ok := fields[field]; ok {
			return value, true
		}
	}
	return "", false
}

// HGetAll gets all the value of a hash key
func (h *HashMap) HGetAll(key string) (map[string]string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	m, ok := h.data[key]
	return m, ok
}

// HDel deletes the entire key or specific member
func (h *HashMap) HDel(key string, mem ...string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(mem) == 0 {
		delete(h.data, key)
	} else {
		for i := 0; i < len(mem); i++ {
			delete(h.data[key], mem[i])
		}
	}
	return true
}
