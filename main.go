package inmemorydb

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Storage struct {
	mu          sync.RWMutex
	data        map[string]interface{}
	flushPath   string
	flushPeriod time.Duration
	stopChan    chan struct{}
}

func NewStorage(flushPath string, flushPeriod time.Duration) *Storage {
	s := &Storage{
		data:        make(map[string]interface{}),
		flushPath:   flushPath,
		flushPeriod: flushPeriod,
		stopChan:    make(chan struct{}),
	}
	s.loadFromDisk()
	go s.flushPeriodically()
	return s
}

func (s *Storage) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Storage) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *Storage) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

func (s *Storage) flushPeriodically() {
	ticker := time.NewTicker(s.flushPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = s.flushToDisk() // Ignore errors for now
		case <-s.stopChan:
			return
		}
	}
}

func (s *Storage) flushToDisk() error {
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

func (s *Storage) loadFromDisk() error {
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

func (s *Storage) Close() error {
	close(s.stopChan)
	return s.flushToDisk()
}
