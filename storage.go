package limiter

import (
	"sync"
	"time"
)

type InMemoryStorage struct {
	data map[string]int
	mx   sync.Mutex
}

// NewInMemoryStorage creates a new InMemoryStorage instance.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]int),
	}
}

// Get retrieves the value associated with the given key from in-memory storage.
func (ims *InMemoryStorage) Get(key string) (int, error) {
	ims.mx.Lock()
	defer ims.mx.Unlock()
	value, ok := ims.data[key]
	if !ok {
		return 0, nil
	}
	return value, nil
}

// Set sets the value associated with the given key in in-memory storage with expiration.
func (ims *InMemoryStorage) Set(key string, value int, expiration time.Duration) error {
	ims.mx.Lock()
	defer ims.mx.Unlock()
	ims.data[key] = value
	go func() {
		time.Sleep(expiration)
		ims.mx.Lock()
		defer ims.mx.Unlock()
		delete(ims.data, key)
	}()
	return nil
}
