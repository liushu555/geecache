package lru

import (
	"fmt"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0))
	lru.Add("key1", String("1234"))
	v, ok := lru.Get("key1")
	fmt.Printf("%T\n", v)
	if !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	fmt.Println(lru.Get("key2"))
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

// 测试是否会自动移除
func TestAdd(t *testing.T) {
	k1, k2, k3 := "key1", "kyey2", "key3"
	v1, v2, v3 := "v1", "v2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap))
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	fmt.Println(lru.Get(k1))
	if _, ok := lru.Get(k1); ok {
		t.Fatalf("cache miss key2 failed")
	}
}
