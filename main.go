package main

import (
	"erinys/erinys"
	"erinys/lru"
	"fmt"
)

// 每个缓存服务都都需要知道所有其它缓存服务节点信息
// 每个缓存服务都要创建一个对应的 group
func startCacheServer(host string, otherpeers ...string) {
	engine := erinys.NewHTTPEngine("/cache", 5, nil)
	engine.RegisterPeer(otherpeers...)
	createGroup(engine, host)
	engine.Run(host)
}

// 每个节点都有自己的 group
func createGroup(engine *erinys.HTTPEngine, peer string) {
	erinys.NewGroup(
		"test",
		erinys.GetterFunc(
			func(key string) (lru.Value, error) {
				m := map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				}
				v, ok := m[key]
				if !ok {
					return nil, fmt.Errorf("%s not exit", key)
				}
				vv := erinys.String(v)
				return vv, nil
			}),
		2<<3,
		engine,
		peer,
	)
}

func main() {

	// 4 个缓存服务的节点
	// 需要有一个注册入口
	otherPeers := []string{
		"127.0.0.1:8001",
		"127.0.0.1:8002",
		"127.0.0.1:8003",
		"127.0.0.1:8004",
	}
	for _, host := range otherPeers {
		go startCacheServer(host, otherPeers...)
	}

	// 对外提供服务
	engine := erinys.NewHTTPEngine("/cache", 5, nil)
	engine.Run(":9001")
}
