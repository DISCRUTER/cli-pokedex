package pokecache

import( 
	"time"
	"sync"
)

var cacheDuration int

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	quit    chan bool
	lock    sync.RWMutex
}


// Set Cache duration
func SetCacheDuration(duration int) {
	cacheDuration = duration
}


// Initialization
func NewCache() *Cache {
	cache := &Cache{
		entries: make(map[string]cacheEntry),
		quit: make(chan bool),
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
		val: value,
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
	ticker := time.NewTicker(time.Duration(cacheDuration) * time.Second)
	defer ticker.Stop()
	lastRefresh := time.Now()
	
	for {
		select {
		case <- c.quit:
			return
		case t := <- ticker.C:
			// Locking resource
			c.lock.Lock()
			// Removing outdated keys
			for key, value := range c.entries {
				if value.createdAt.Before(lastRefresh) {
					delete(c.entries, key)
					continue
				}
			}
			// Unlocking Resource
			c.lock.Unlock()
			// updating lastRefresh
			lastRefresh = t
		}
	}
}
