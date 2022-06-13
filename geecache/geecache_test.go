package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var g Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, _ := g.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

//模拟数据库
var db = map[string]string{
	"test1": "90",
	"test2": "80",
	"test3": "70",
}

func TestGet(t *testing.T) {
	//统计回调次数
	loadCounts := make(map[string]int, len(db))
	geeCache := NewGroup("scores", 512, GetterFunc(
		func(key string) ([]byte, error) {
			if v, ok := db[key]; ok {
				log.Println("searched!")
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s is not existed", key)
		}))

	for k, v := range db {
		if view, err := geeCache.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value by callback function")
		}
		if _, err := geeCache.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := geeCache.Get("unknown"); err == nil {
		t.Fatalf("unkown doesn't exist, but got %s", view)
	}
}
