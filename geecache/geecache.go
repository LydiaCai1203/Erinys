package geecache

import (
	"fmt"
	"log"
	"myRPC/geecache/singleflight"
	"sync"
)

// 当用户获取数据时发现数据不存在
// 则触发回调(具体由用户实现)，得到源数据
type Getter interface {
	Get(key string) ([]byte, error)
}

// 定义一个函数类型 F，实现接口 A 的方法，然后在方法中调用自己
// 是 Go 中将其它函数转换为接口 A 的常用技巧
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 每个 Group 可以认为是一个缓存的命名空间
type Group struct {
	name      string
	getter    Getter // 缓存未命中时获取源数据的回调
	mainCache cache  // 一开始实现的并发缓存
	peers     PeePicker
	loader    *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()         // 写锁获取
	defer mu.Unlock() // 写锁释放
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// 加入缓存
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	// 如果缓存查不到就
	return g.load(key)
}

func (g *Group) RegisterPeers(peers PeePicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) getFromPeer(peer PeeGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(
		key,
		func() (interface{}, error) {
			if g.peers != nil {
				if peer, ok := g.peers.PickPeer(key); ok {
					if value, err = g.getFromPeer(peer, key); err == nil {
						return value, nil
					}
					log.Println("[GeeCache] Failed to get from peer", err)
				}
			}
			return g.getLocally(key)
		})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}
