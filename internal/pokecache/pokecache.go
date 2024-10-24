package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]CacheItem
	mux sync.RWMutex
}

type CacheItem struct {
	val []byte
	createdAt time.Time
}

func (c *Cache) Set(key string, val []byte) error {
	c.mux.Lock()
	if len(c.cacheMap[key].val) != 0 {
		return fmt.Errorf("key %s already exists", key)
	}
	c.mux.Unlock()

	item := CacheItem{
		val: val,
        createdAt: time.Now(),
	}

	c.mux.Lock()
	c.cacheMap[key] = item
	c.mux.Unlock()

	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
    item, ok := c.cacheMap[key]
	if !ok {
        return []byte{}, false
    }
    return item.val, ok
}

func (c *Cache) Delete(key string) {
    delete(c.cacheMap, key)
}

func NewCache() *Cache {
	return &Cache{
        cacheMap: make(map[string]CacheItem),
		mux: sync.RWMutex{},
    }
}