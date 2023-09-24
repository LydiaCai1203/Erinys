package erinys

import (
	"erinys/lru"
	"fmt"
	"sync"
)

type Getter interface {
	// 从源站获取数据
	Get(key string) (lru.Value, error)
}

type GetterFunc func(key string) (lru.Value, error)

func (f GetterFunc) Get(key string) (lru.Value, error) {
	v, err := f(key)
	return v, err
}

type Group struct {
	name   string
	getter Getter
	cache  *lru.SafeLruCache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, getter Getter, capacity int64) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		cache:  lru.NewLruCache(capacity),
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := groups[name]
	return g
}

// 尝试从本地缓存获取数据，获取不到则从源站获取
func (g *Group) Get(key string) (lru.Value, error) {
	// 万一没有给 getter 赋值
	if g.getter == nil {
		return nil, fmt.Errorf("getter not assigned")
	}
	// 看本地缓存有无数据
	v, err := g.cache.Get(key)
	if err == nil {
		fmt.Println("hit cache")
		return v, nil
	}
	// 没有数据则去原站数据访问
	v, _ = g.getter.Get(key)
	g.cache.Add(key, v)
	return v, nil
}
