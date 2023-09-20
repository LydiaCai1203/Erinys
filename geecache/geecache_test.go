package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGetter(t *testing.T) {
	// 相当于定义了一个这个函数, 这个函数还有一个 GET 方法
	// func f(key string) ([]byte, error) {...}
	// 这样做的好处是避免通过定义 struct 然后再定义方法，可以少写几行，大概这意思
	var f Getter = GetterFunc(
		func(key string) ([]byte, error) {
			v := db[key]
			return []byte(v), nil
		},
	)
	expect := []byte("589")
	// 用于比较两个切片是否相等
	if v, _ := f.Get("Jack"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

func TestGet(t *testing.T) {
	// 自定义了 cache miss 时的回调函数
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup(
		"scores",
		2<<10,
		GetterFunc(
			func(key string) ([]byte, error) {
				log.Println("[SlowDB] search key", key)
				if v, ok := db[key]; ok {
					if _, ok := loadCounts[key]; !ok {
						loadCounts[key] = 0
					}
					loadCounts[key] += 1
					return []byte(v), nil
				}
				return nil, fmt.Errorf("%s not exist", key)
			}))

	// 依次访问 db 里的每个 key, 如果成功获取了源数据并缓存了，loadCount[k] 为 1
	// 如果每次都从源数据获取说明没有缓存成功，则 loadCounts[k] > 1
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
