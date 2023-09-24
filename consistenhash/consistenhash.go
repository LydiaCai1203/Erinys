package consistenhash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func([]byte) uint32

type PeerPool struct {
	hash     HashFunc
	peermap  map[int]string // key: 虚拟节点，value: 真实节点
	replicas int            // 每个真实节点有的副本节点数量
	vpeers   []int          // 虚拟节点环
}

func NewPeerPool(replicas int, fn HashFunc) *PeerPool {
	pool := &PeerPool{
		hash:     fn,
		replicas: replicas,
		vpeers:   make([]int, 0),
		peermap:  make(map[int]string),
	}
	if pool.hash == nil {
		pool.hash = crc32.ChecksumIEEE
	}
	return pool
}

// 节点注册
func (pool *PeerPool) RegisterPeer(peers ...string) {
	// 真实节点
	for _, peer := range peers {
		// 副本节点
		for i := 0; i < pool.replicas; i++ {
			vpeer := peer + "-" + strconv.Itoa(i)
			vpeerBytes := []byte(vpeer)
			hash := int(pool.hash(vpeerBytes))
			pool.vpeers = append(pool.vpeers, hash)
			pool.peermap[hash] = peer
		}
	}
	sort.Ints(pool.vpeers)
}

// 输入数据 key 查询 key 应当存储的真实节点名字
func (pool *PeerPool) GetPeerByKey(key string) (peer string) {
	keyByte := []byte(key)
	hash := int(pool.hash(keyByte))
	idx := sort.Search(
		len(pool.vpeers),
		func(i int) bool {
			return pool.vpeers[i] >= hash
		})
	idx = idx % len(pool.vpeers)
	vhash := pool.vpeers[idx]
	return pool.peermap[vhash]
}
