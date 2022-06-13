package main

import (
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

func main() {
	geecache.NewGroup("scores", 512, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDB] search key:" + key)
			val, err := db[key]
			if !err {
				return nil, fmt.Errorf("%s not existed", key)
			}
			return []byte(val), nil
		},
	))

	address := "127.0.0.1:9999"
	httpPool := geecache.NewHTTPPool(address)

	log.Println("geecache is running at:", address)
	log.Fatal(http.ListenAndServe(address, httpPool))
}
