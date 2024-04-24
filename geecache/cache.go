package geecache

import (
	"geecache/lru"
	"sync"
)

//并发控制

type cache struct {
	//并发控制锁
	mu sync.Mutex
	//实现lru
	lru *lru.Iru
	//最大存储值
	maxBytes int64
}

func (c *cache) Add(key string, value Byteview) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//延迟初始化
	if c.lru == nil {
		c.lru = lru.New(c.maxBytes)
	}
	c.lru.Add(key, value)
}

func (c *cache) Get(key string) (value Byteview, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(Byteview), ok
	}
	return
}
