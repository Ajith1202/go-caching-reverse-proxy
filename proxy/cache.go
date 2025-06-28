package proxy

import (
	"sync"
	"time"
)

type CacheItem struct {
	value  interface{}
	expiry time.Time
}

type Cache struct {
	data map[string]CacheItem
	mu   sync.RWMutex
}

func NewCache() *Cache {
	c := &Cache{
		data: make(map[string]CacheItem),
	}
	return c
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheItem{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, exists := c.data[key]
	c.mu.RUnlock()

	if !exists || time.Now().After(item.expiry) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return nil, false
	}

	return item.value, exists
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()
			c.mu.Lock()
			for k, v := range c.data {
				if now.After(v.expiry) {
					delete(c.data, k)
				}
			}
			c.mu.Unlock()
		}
	}()
}
