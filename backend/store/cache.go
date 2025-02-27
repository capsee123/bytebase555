package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

// cacheNamespace is the type of a cache.
type cacheNamespace string

const (
	// pipelineCacheNamespace is the cache type of pipelines.
	pipelineCacheNamespace cacheNamespace = "pl"
)

// CacheService implements a cache.
type CacheService struct {
	sync.Mutex
	cache map[string][]byte
}

// newCacheService creates a cache service.
func newCacheService() *CacheService {
	return &CacheService{
		cache: make(map[string][]byte),
	}
}

// FindCache finds the value in cache.
func (s *CacheService) FindCache(namespace cacheNamespace, id int, entry interface{}) (bool, error) {
	key := generateKey(namespace, id)

	s.Lock()
	defer s.Unlock()
	value, exists := s.cache[key]
	if exists {
		dec := gob.NewDecoder(bytes.NewReader(value))
		if err := dec.Decode(entry); err != nil {
			return false, errors.Wrapf(err, "failed to decode entry for cache namespace: %s", namespace)
		}
		return true, nil
	}
	return false, nil
}

// UpsertCache upserts the value to cache.
func (s *CacheService) UpsertCache(namespace cacheNamespace, id int, entry interface{}) error {
	key := generateKey(namespace, id)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(entry); err != nil {
		return errors.Wrapf(err, "failed to encode entry for cache namespace: %s", namespace)
	}

	s.Lock()
	defer s.Unlock()
	s.cache[key] = buf.Bytes()

	return nil
}

// DeleteCache deletes the key from cache.
func (s *CacheService) DeleteCache(namespace cacheNamespace, id int) {
	key := generateKey(namespace, id)

	s.Lock()
	defer s.Unlock()
	delete(s.cache, key)
}

func generateKey(namespace cacheNamespace, id int) string {
	return fmt.Sprintf("%s%d", namespace, id)
}
