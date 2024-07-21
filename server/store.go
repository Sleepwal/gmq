package server

import (
	"fmt"
	"sync"
)

type StoreProducerFunc func() Storage

type Storage interface {
	Push([]byte) (int, error)
	Get(offset int) ([]byte, error)
}

type MemoryStore struct {
	mu   sync.RWMutex
	data [][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make([][]byte, 0),
	}
}

func (m *MemoryStore) Push(bytes []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = append(m.data, bytes)
	return len(m.data) - 1, nil
}

func (m *MemoryStore) Get(offset int) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be smaller than 0")
	}
	if offset > len(m.data)-1 {
		return nil, fmt.Errorf("offset cannot be bigger than %d", len(m.data)-1)
	}
	return m.data[offset], nil
}
