package storage

import (
	"sync"
)

type InMemoryStorage struct {
	mu    sync.RWMutex
	cache map[string]*Bucket
}

func NewInMemory() *InMemoryStorage {
	return &InMemoryStorage{
		mu:    sync.RWMutex{},
		cache: make(map[string]*Bucket),
	}
}

func (s *InMemoryStorage) GetBucket(key, namespace string) *Bucket {
	compositeKey := key + ":" + namespace

	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.cache[compositeKey]
	if !ok {
		b = NewBucket(5.0, 0.5)
		s.cache[compositeKey] = b
	}

	return b
}
