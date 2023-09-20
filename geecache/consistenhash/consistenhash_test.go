package consistenhash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	// 这就依赖注入了？
	hash := New(
		3,
		func(key []byte) uint32 {
			i, _ := strconv.Atoi(string(key))
			return uint32(i)
		},
	)

	hash.Add("6", "4", "2")
	testCases := map[string]string{
		"2":  "2", // key=2 的数据应当存在节点 2 上
		"11": "2", // key=11 的数据应当存在节点 2 上
		"23": "4", // key=23 的数据应当存在节点 4 上
		"27": "2", // key=27 的数据应当存在节点 2 上
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	hash.Add("8")
	testCases["27"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
