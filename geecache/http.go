package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geeCache/"

type HTTPPool struct {
	//主机地址以及端口
	self string
	//基本路径
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (h *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("server %s : %s", h.self, fmt.Sprintf(format, v...))
}

//约定访问路径为  /<basepath>/<groupname>/<key>
//实现serveHTTP即实现了HTTP的Handler接口
func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("unexpected path:" + r.URL.Path)
	}

	h.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad Request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	keyName := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, fmt.Sprintf("group %s not found", groupName), http.StatusNotFound)
	}
	value, err := group.Get(keyName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.ByteSlice())

}
