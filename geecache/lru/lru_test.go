/*
cd myRPC/geecache/lru
go test -run TestGet

https://geektutu.com/post/quick-go-test.html
*/

package lru

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveTail(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)

	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		// 测试出错则输出错误信息并中止
		t.Fatalf("RemoveTail key1 failed")
	}
}

func TestMetux(t *testing.T) {
	var m sync.Mutex
	var set = make(map[int]bool, 0)

	f := func(num int) {
		// 保证一个时刻只有一个协程能访问 set
		m.Lock()
		defer m.Unlock()
		if _, exist := set[num]; !exist {
			fmt.Println(num)
		}
		set[num] = true
	}

	for i := 0; i < 10; i++ {
		go f(100)
	}

	time.Sleep(time.Second)
}
