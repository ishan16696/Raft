package db

import (
	"sync"
)

// The key-value database.
type DB struct {
	data  map[string]interface{}
	mutex sync.RWMutex
}

// Creates a new database.
func New() *DB {
	return &DB{
		data: make(map[string]interface{}),
	}
}

// Retrieves the value for a given key.
func (db *DB) Get(key string) interface{} {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	return db.data[key]
}

// Sets the value for a given key.
func (db *DB) Put(key string, value interface{}) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.data[key] = value
}