package erinys

import (
	"erinys/lru"
	"fmt"
	"testing"
)

func testget(key string) (lru.Value, error) {
	m := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	v, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("%s not exit", key)
	}
	vv := String(v)
	return vv, nil
}

func TestGroupGet(t *testing.T) {
	group := NewGroup(
		"test",
		GetterFunc(testget),
		2<<3,
	)
	group.Get("key1")
	group.Get("key1")
}
