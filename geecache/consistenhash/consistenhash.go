package consistenhash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 函数类型
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	// 默认 hash 函数
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 添加虚拟节点
// keys: 一个或多个真实节点的名字
// m.replicas: 一个或多个虚拟节点的名字
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 把虚拟节点全都加环上了
			m.keys = append(m.keys, hash)
			// 存储的是虚拟节点和真实节点的映射关系
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// 查询对应的数据应当存储的真实节点信息
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	// 数据 hash
	hash := int(m.hash([]byte(key)))
	// 找到虚拟节点的 hash 值
	idx := sort.Search(
		len(m.keys),
		func(i int) bool {
			return m.keys[i] >= hash
		},
	)
	// 如果 idx = len(m.keys), 说明数据其实在第一个节点上
	// 因为 m.keys 其实是一个环
	idxKey := m.keys[idx%len(m.keys)]
	// 找到真实节点的 hash 值
	return m.hashMap[idxKey]
}
