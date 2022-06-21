package main

import (
	"flag"
	"fmt"
	"geeCache/geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *geecache.Group {
	return geecache.NewGroup("scores", 512, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDB] search key:" + key)
			val, err := db[key]
			if !err {
				return nil, fmt.Errorf("%s not existed", key)
			}
			return []byte(val), nil
		},
	))
}

//启动服务
func startCacheServer(addr string, address []string, gee *geecache.Group) {
	peers := geecache.NewHTTPPool(addr)
	peers.Set(address...)
	gee.RegisterPeers(peers)
	log.Println("server is running at:", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

//开启API服务
func startAPIServer(apiAddr string, gee *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		},
	))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}
