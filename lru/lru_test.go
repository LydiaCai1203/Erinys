package lru

import (
	"fmt"
	"testing"
)

type mystr string

func (s mystr) Len() int64 {
	return int64(len(s))
}

func TestGet(t *testing.T) {
	l := New(2 << 10)
	var v mystr = "caiqj"
	l.Add("name", v)
	value, err := l.Get("name")
	if err != nil {
		t.Fatalf("Get key {name} Failed")
	}
	fmt.Printf("Get key {name} Success: %v", value)
}

func TestREmoveTail(t *testing.T) {
	l := New(2 << 3)
	l.Add("name", mystr("caiqj"))
	l.Add("age", mystr("22"))
	l.Add("from", mystr("zhejiang"))
	l.All()
}
