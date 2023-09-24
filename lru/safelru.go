package lru

import (
	"sync"
)

// 协程安全结构
// 可以同时 Get, Add 时不能 Add/Get
type SafeLruCache struct {
	mu  sync.Mutex
	lru *LruCache
}

func NewLruCache(capacity int64) *SafeLruCache {
	return &SafeLruCache{
		lru: New(capacity),
	}
}

func (l *SafeLruCache) Get(key string) (value Value, err error) {
	// 读锁
	l.mu.Lock()
	v, err := l.lru.Get(key)
	l.mu.Unlock()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (l *SafeLruCache) Add(key string, value Value) {
	// 写锁
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lru.Add(key, value)
}
