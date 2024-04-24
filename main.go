package main

import (
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var getter1 geecache.Getter = geecache.GetterFunc(func(key string) ([]byte, error) {
	log.Println("[SlowDB] search key", key)
	if v, ok := db1[key]; ok {
		return []byte(v), nil
	}
	return []byte{}, fmt.Errorf("%s not exist", key)
})

var db1 = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// 测试HTTP服务
func main() {
	geecache.NewGroup("test", getter1, 1000)
	server := geecache.NewHTTPool("localhost:9000")
	fmt.Println("服务启动")
	http.ListenAndServe("localhost:9000", server)
	//log.Fatal(http.ListenAndServe("localhost:9000", server))
}
