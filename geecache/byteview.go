package geecache

type ByteView struct {
	b []byte // 存储真实的缓存值
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// 被缓存的对象一定要实现 Value interface
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 只读，因此返回拷贝防止被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}
