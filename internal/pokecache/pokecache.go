package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	quit     chan bool
	lock     sync.RWMutex
	interval time.Duration
}

// Initialization
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		quit:     make(chan bool),
		interval: interval,
	}

	// Intializing cleanup
	go cache.reapLoop()

	return cache
}

// Adding new entry
func (c *Cache) Add(key string, value []byte) {
	// Locking the resource
	c.lock.Lock()
	defer c.lock.Unlock()

	// Adding value to entries
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

// Getting value
func (c *Cache) Get(key string) ([]byte, bool) {
	// Locking the resource
	c.lock.RLock()
	defer c.lock.RUnlock()

	// Fetching the value
	value, exist := c.entries[key]
	if !exist {
		return nil, false
	}
	return value.val, true
}

func (c *Cache) Stop() {
	c.quit <- true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.quit:
			return
		case <-ticker.C:
			// Locking resource
			c.lock.Lock()
			// Removing outdated keys
			for key, value := range c.entries {
				if time.Since(value.createdAt) >= c.interval {
					delete(c.entries, key)
				}
			}
			// Unlocking Resource
			c.lock.Unlock()
		}
	}
}
