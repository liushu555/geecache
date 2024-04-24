package geecache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// 接口型函数
type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

type Group struct {
	//该缓存属于哪一组
	name string
	//缓存未命中时的回调函数
	getter Getter
	//该组数据的存储地
	mainCache cache
}

var (
	groups = make(map[string]*Group)
	mu     sync.RWMutex
)

func NewGroup(name string, getter Getter, maxBytes int64) *Group {
	if getter == nil {
		panic("Getter nil")
	}
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{maxBytes: maxBytes},
	}
	mu.Lock()
	defer mu.Unlock()
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	if g, ok := groups[name]; ok {
		return g
	}
	defer mu.RUnlock()
	return nil
}

func (g *Group) Get(key string) (Byteview, error) {
	if key == "" {
		return Byteview{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[GeeCache] hit")
		//缓存命中
		return v, nil
	}
	//缓存未命中，调用回调，加载数据

	return g.load(key)
}

func (g *Group) load(key string) (Byteview, error) {
	v, err := g.getter.Get(key)
	if err != nil {
		return Byteview{}, err
	}
	b := Byteview{bytes: cloneByte(v)}
	//将回调得到的数值写入内存
	g.mainCache.Add(key, b)
	return b, nil
}
