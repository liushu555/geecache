package geecache

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

var getter Getter = GetterFunc(func(key string) ([]byte, error) {
	log.Println("[SlowDB] search key", key)
	if v, ok := db[key]; ok {
		return []byte(v), nil
	}
	return []byte{}, fmt.Errorf("%s not exist", key)
})

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// 测试getter函数接口
func TestGetter(t *testing.T) {
	log.Println("teststart ")
	gee := NewGroup("test", getter, 1000)
	log.Println("test init ")
	for k, v := range db {
		log.Println("range db", k, v)
		if b, err := gee.Get(k); err != nil || b.String() != v {
			t.Fatalf("failed to get value of %s", k)
		}
	}
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

// 测试HTTP服务
func TestHttp(t *testing.T) {
	NewGroup("test", getter, 1000)
	server := NewHTTPool("localhost:9000")
	http.ListenAndServe("localhost:9000", server)
	//log.Fatal(http.ListenAndServe("localhost:9000", server))
}
