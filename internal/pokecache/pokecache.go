package pokecache

import (
	"time"
)

type Cache struct {
	cacheMap map[string]CacheItem
}

type CacheItem struct {
	val       []byte
	createdAt time.Time
}

func (c *Cache) Add(key string, val []byte) {
	c.cacheMap[key] = CacheItem{
		val:       val,
		createdAt: time.Now().UTC(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	item, ok := c.cacheMap[key]
	return item.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}


func (c *Cache) reap(interval time.Duration) {
	timeToDelete := time.Now().Add(-interval)
	for k,v := range c.cacheMap {
		if v.createdAt.Before(timeToDelete) {
			delete(c.cacheMap, k)
		}
	}
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cacheMap: make(map[string]CacheItem),
	}
	go c.reapLoop(interval)
	return c
}
