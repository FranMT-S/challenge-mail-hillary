package models

import (
	"sort"
	"strconv"
	"sync"
)

// SafeMap con RWMutex
type SafeMap struct {
	mu sync.RWMutex
	m  map[string]PageResult
}

// NewSafeMap creates a new SafeMap
func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]PageResult),
	}
}

// Set sets a value
func (s *SafeMap) Set(key string, value PageResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func (s *SafeMap) SetMap(m map[string]PageResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = m
}

func (s *SafeMap) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}

// Get gets a value
func (s *SafeMap) Get(key string) (PageResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.m[key]
	return val, ok
}

// Delete deletes a value
func (s *SafeMap) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
}

func (s *SafeMap) GetCopy() map[string]PageResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copyMap := make(map[string]PageResult, len(s.m))
	for k, v := range s.m {
		copyMap[k] = v
	}

	return copyMap
}

// Range iterates over all elements
func (s *SafeMap) Range(f func(key int, value PageResult)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]int, 0, len(s.m))

	for _, v := range s.m {
		keys = append(keys, v.Page)
	}

	sort.Ints(keys)

	for _, k := range keys {
		f(k, s.m[strconv.Itoa(k)])
	}
}
