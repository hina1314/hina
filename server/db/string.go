package db

import (
	"sync"
)

type String struct {
	Data map[string]string
	mu   sync.RWMutex
}

func NewStrings() *String {
	return &String{
		Data: make(map[string]string),
	}
}

func (db *String) Set(key, value string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Data[key] = value
	return true
}

func (db *String) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.Data[key]
	return val, ok
}

func (db *String) GetAll() map[string]string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.Data
}

func (db *String) Del(key string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.Data, key)
	return true
}
