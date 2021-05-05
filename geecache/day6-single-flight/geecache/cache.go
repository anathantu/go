package geecache

import (
	"example1/geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		//这里的v.(ByteView)属于接口转换成类型，
		//基本格式是：t := i.(T)
		//i 代表接口变量，T 代表转换的目标类型，t 代表转换后的变量。
		return v.(ByteView), ok
	}

	return
}
