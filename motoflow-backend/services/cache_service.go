package services

import (
	"sync"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration int64
}

// MemoryCache provides in-memory thread-safe key-value caching with TTL
type MemoryCache struct {
	items map[string]cacheItem
	mu    sync.RWMutex
}

var GlobalCache *MemoryCache
var cacheOnce sync.Once

func GetCache() *MemoryCache {
	cacheOnce.Do(func() {
		GlobalCache = &MemoryCache{
			items: make(map[string]cacheItem),
		}
		// Background cleanup of expired items every minute
		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			for range ticker.C {
				GlobalCache.cleanExpired()
			}
		}()
	})
	return GlobalCache
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	c.items[key] = cacheItem{
		value:      value,
		expiration: exp,
	}
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		return nil, false
	}

	return item.value, true
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *MemoryCache) InvalidatePrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	prefixLen := len(prefix)
	for k := range c.items {
		if len(k) >= prefixLen && k[:prefixLen] == prefix {
			delete(c.items, k)
		}
	}
}

func (c *MemoryCache) cleanExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for k, item := range c.items {
		if item.expiration > 0 && now > item.expiration {
			delete(c.items, k)
		}
	}
}

