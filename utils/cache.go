package utils

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

var cache = make(map[string]CacheEntry)
var mu sync.RWMutex

func SetCache(key string, data interface{}, duration time.Duration) {
	mu.Lock()
	defer mu.Unlock()
	cache[key] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(duration),
	}
}

func GetCache(key string) (interface{}, bool) {
	mu.RLock()
	defer mu.RUnlock()

	entry, exists := cache[key]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Data, true
}
