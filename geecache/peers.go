package geecache

type PeePicker interface {
	// 根据传入的 key 选择相应节点 PeerGetter
	PickPeer(key string) (peer PeeGetter, ok bool)
}

type PeeGetter interface {
	// 用于从对应的 group 中查找缓存值
	Get(group string, key string) ([]byte, error)
}
