package lru

import (
	"container/list"
)

type Cache struct {

	//存储k,v实现lru
	ll *list.List
	//k-k,v-ll中的一个元素
	cache map[string]*list.Element
	//最大存储字节数
	maxBytes int64
	//当前存储的字节数
	curBytes int64

	//删除回调函数
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// 初始化
func New(maxBytes int64) *Cache {
	return &Cache{
		ll:       &list.List{},
		cache:    make(map[string]*list.Element),
		maxBytes: maxBytes,
		curBytes: 0,
	}
}

// 增加或者更新
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		//更新
		c.ll.MoveToFront(ele)
		//强转类型
		kv := ele.Value.(*entry)
		//更新字节数
		c.curBytes = c.curBytes - int64(kv.value.Len()) + int64(value.Len())
		//更新值
		kv.value = value
	} else {
		//新增
		ele := c.ll.PushFront(&entry{key, value})
		//插入cache,做绑定
		c.cache[key] = ele
		//更新字节数
		c.curBytes += int64(len(key)) + int64(value.Len())
	}
	//判断是否超出最大内存，超出则进行删除
	for c.maxBytes != 0 && c.curBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// 删除最近未被访问的元素
func (c *Cache) RemoveOldest() {
	//删除队尾元素
	ele := c.ll.Back()
	if ele != nil {
		//更新链表
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		//更新字节
		c.curBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//更新map
		delete(c.cache, kv.key)
	}
}

// 查询
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		//存在
		kv := ele.Value.(*entry)
		//更新lru
		c.ll.MoveToFront(ele)
		return kv.value, true
	}
	return
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
