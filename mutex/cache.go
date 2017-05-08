package main

import (
	"fmt"
	"sync"
	"time"
)

//CacheEntry represents an entry in the cache
type CacheEntry struct {
	value   string
	expires time.Time
}

//Cache represents a map[string]string that is safe
//for concurrent access
type Cache struct {
	//protect this map with a RWMutex
	mu      sync.RWMutex // to create concurrent safe cache!!!
	entries map[string]*CacheEntry
	quit    chan bool
}

//NewCache creates and returns a new Cache
func NewCache() *Cache {
	c := &Cache{
		entries: make(map[string]*CacheEntry),
		mu:      sync.RWMutex{},
		quit:    make(chan bool),
	}
	// start janitor here
	go c.startJanitor() // runs independently since "go"
	return c
}

func (c *Cache) Close() {
	c.quit <- true
}

func (c *Cache) startJanitor() { // cache will be this function's receiver
	ticker := time.NewTicker(time.Second)
	// ticker will be a channel, time package will write something to this channel every second
	for { // infinate for-loop
		select { // non-blocking channel
		case <-ticker.C:
			c.purgeExpired()
		case <-c.quit:
			return // to be garbage collected
		}
	}
}

func (c *Cache) purgeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	nPurged := 0
	for key, entry := range c.entries {
		if now.After(entry.expires) { // if now is after entry.expires
			delete(c.entries, key)
			nPurged++
		}
	}
	fmt.Printf("Purged %d entries\n", nPurged)
}

//Get returns the value associated with the requested key.
//The returned boolean will be false if the key was not
//in the cache.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock() // ALWAYS follow RLock()
	//implement this method and
	//replace the return statement below
	entry := c.entries[key]
	if entry == nil {
		return "", false
	}
	return entry.value, true
}

//Set sets the value associated with the given key.
//If the key is not yet in the cache, it will be added.
func (c *Cache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//implement this method
	entry := c.entries[key]
	if entry == nil {
		entry = &CacheEntry{}
		c.entries[key] = entry
	}
	entry.value = value
	entry.expires = time.Now().Add(ttl) // ttl = time to live
	//(duration of how long client wants this to live in the cache)
}
