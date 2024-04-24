package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

//运行服务端，方便其他节点进行数据访问

const basePath = "/_geecache/"

type HTTPPool struct {
	//addr localhost:9000
	self string

	//basePath
	basePath string
}

func NewHTTPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: basePath,
	}
}

func (p *HTTPPool) Log(method string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(method, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//检查路径是否符合规则/{basePath}/{groupname}/{key}
	fmt.Println(r.URL.Path, p.basePath)
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("this request not has the right prefix")
	}
	seq := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(seq) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	//获取groupname以及key，返回查找到的值，如何去查找？
	groupname := seq[0]
	key := seq[1]

	g := GetGroup(groupname)
	if g == nil {
		http.Error(w, "group is not found", http.StatusNotFound)
		return
	}
	view, err := g.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
