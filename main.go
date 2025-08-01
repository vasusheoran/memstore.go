package inmemorydb

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Storage interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
	All() map[string]interface{}
	Close() error
	Flush() error
}

type storage struct {
	mu          sync.RWMutex
	data        map[string]interface{}
	flushPath   string
	flushPeriod time.Duration
	stopChan    chan struct{}
}

func NewStorage(flushPath string, flushPeriod time.Duration) Storage {
	s := &storage{
		data:        make(map[string]interface{}),
		flushPath:   flushPath,
		flushPeriod: flushPeriod,
		stopChan:    make(chan struct{}),
	}
	s.loadFromDisk()

	if flushPeriod != 0 {
		go s.flushPeriodically()
	}

	return s
}

func (s *storage) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *storage) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *storage) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// All returns a copy of the underlying map for read-only purposes
func (s *storage) All() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a shallow copy to avoid race conditions
	copy := make(map[string]interface{}, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}

func (s *storage) flushPeriodically() {
	ticker := time.NewTicker(s.flushPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = s.Flush() // Ignore errors for now
		case <-s.stopChan:
			return
		}
	}
}

func (s *storage) Flush() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Create(s.flushPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(s.data)
}

func (s *storage) loadFromDisk() error {
	file, err := os.Open(s.flushPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&s.data)
}

func (s *storage) Close() error {
	close(s.stopChan)
	return s.Flush()
}
