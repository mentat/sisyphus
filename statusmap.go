package main

import "sync"

// Status ..
type Status struct {
	Running bool   `json:"running"`
	Done    bool   `json:"done"`
	Error   string `json:"error,omitempty"`
}

// ThreadSafeMap ...
type ThreadSafeMap struct {
	data  map[string]Status
	mutex sync.RWMutex
}

// NewThreadSafeMap ...
func NewThreadSafeMap() *ThreadSafeMap {
	tsm := &ThreadSafeMap{
		data: make(map[string]Status),
	}
	return tsm
}

// Delete ...
func (m *ThreadSafeMap) Delete(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.data, key)
}

// Set ...
func (m *ThreadSafeMap) Set(key string, data Status) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = data
}

// Get ...
func (m *ThreadSafeMap) Get(key string) (data Status, ok bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	data, ok = m.data[key]
	return
}
