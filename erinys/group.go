package erinys

import (
	"erinys/lru"
	"fmt"
	"strings"
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
	peer   PeerPicker
	self   string // group 所在节点
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(
	name string,
	getter Getter,
	capacity int64,
	peer PeerPicker,
	self string,
) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		cache:  lru.NewLruCache(capacity),
		peer:   peer,
		self:   self,
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

// 根据 key 获取数据
func (g *Group) Get(key string) (lru.Value, error) {
	// 需要根据 key 在服务里找到 peer 信息
	_, peer := g.peer.PickPeer(key)
	fmt.Printf("key: %s, self: %v, dst: %v\n", key, g.self, peer)
	if peer != g.self {
		v, err := g.getFromPeer(key, g.name)
		return v, err
	}
	v, err := g.getFromLocal(key)
	return v, err
}

// 从远程节点获取数据
func (g *Group) getFromPeer(key string, group string) (lru.Value, error) {
	groupname := strings.Split(group, "-")[0]
	peerclient, _ := g.peer.PickPeer(key)
	v, err := peerclient.Get(groupname, key)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// 从本地获取数据
func (g *Group) getFromLocal(key string) (lru.Value, error) {
	// 本地请求
	if g.getter == nil {
		return nil, fmt.Errorf("getter not assigned")
	}
	// 先看本地缓存有无数据
	v, err := g.cache.Get(key)
	if err == nil {
		fmt.Println("hit cache")
		return v, nil
	}
	// 没有数据则去原站数据访问
	v, _ = g.getter.Get(key)
	if v != nil {
		g.cache.Add(key, v)
	}
	return v, nil
}
