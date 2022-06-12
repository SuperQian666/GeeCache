package geeCache

import (
	"geeCache/lru"
	"sync"
)

type cache struct {
	mutex         sync.Mutex
	lru           *lru.Cache
	cacheCapBytes int64
}

func (c *cache) Add(key string, value ByteView) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheCapBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
