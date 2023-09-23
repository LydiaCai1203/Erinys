package lru

import (
	"container/list"
	"fmt"
)

type Value interface {
	Len() int64 // 返回自身占用空间大小
}

type Element struct {
	key   string
	value Value
}

type LruCache struct {
	capacity int64
	cursize  int64
	dict     map[string]*list.Element
	dllist   *list.List
}

func New(capcity int64) *LruCache {
	return &LruCache{
		capacity: capcity,
		dict:     make(map[string]*list.Element),
		dllist:   list.New(),
	}
}

// Add: 存入 key-value 到 LruCache 里
func (l *LruCache) Add(key string, value Value) {
	eptr, ok := l.dict[key]
	if ok {
		// 更新存在节点 value && 容量信息
		l.dllist.MoveToFront(eptr)
		v := eptr.Value.(*Element)
		v.value = value
		diff := value.Len() - v.value.Len()
		l.cursize += diff
	} else {
		// 新插入节点 && 更新容量信息
		ele := &Element{key, value}
		node := l.dllist.PushFront(ele)
		l.dict[key] = node
		diff := value.Len() + int64(len(key))
		l.cursize += diff
	}
	// 这样写仍有问题，应该递归删除队尾元素
	// 如果当前容量超出了限定容量，则删除队尾元素
	if l.cursize > l.capacity {
		l.RemoveTail()
	}
}

// 删除最后一个节点 && 更新容量信息
func (l *LruCache) RemoveTail() {
	node := l.dllist.Back()
	l.dllist.Remove(node)
	ele := node.Value.(*Element)
	nodeSize := ele.value.Len() + int64(len(ele.key))
	l.cursize -= nodeSize
}

// 访问缓存，更新节点在列表中的位置信息
func (l *LruCache) Get(key string) (value Value, err error) {
	eptr, ok := l.dict[key]
	if !ok {
		return nil, fmt.Errorf("key: %s not exist", key)
	}
	l.dllist.MoveToFront(eptr)
	ele := eptr.Value.(*Element)
	return ele.value, nil
}

func (l *LruCache) All() {
	for i := l.dllist.Front(); i != nil; i = i.Next() {
		fmt.Println(i.Value)
	}
}
