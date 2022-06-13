package lru

import (
	"fmt"
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(0, nil)
	testKey := "key1"
	testValue := String("1234")
	lru.Add(testKey, testValue)
	if value, ok := lru.Get("key1"); !ok || string(value.(String)) != string(testValue) {
		t.Fatalf("cache hit %s failed", testKey)
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "v1", "v2", "v3"
	cap := len(k1 + v1 + k2 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("k1"); ok {
		t.Fatalf("false")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	fmt.Println(len(keys))
	fmt.Println(cap(keys))
	callback := func(key string, value Value) {
		keys = append(keys, key)
		fmt.Println(cap(keys))
	}
	lru := New(10, callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
