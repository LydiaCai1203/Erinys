package singleflight

import "sync"

// 代表正在运行或已经结束了的请求
type call struct {
	wg  sync.WaitGroup
	val interface{} // 用户查询的 key 所对应的 val
	err error       // 查询结果存储
}

// 管理不同 key 的请求
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// 针对相同的 key, 无论 Do 被调用多少次, fn 只会被调用一次，fn 调用结束，再返回返回值
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         // 如果请求正在执行中则等待
		return c.val, c.err // 请求结束啧返回结果
	}

	c := new(call)
	g.m[key] = c // 已经有 key 对应的请求正在处理
	g.mu.Unlock()
	c.val, c.err = fn() // 调用 fn 发起请求
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err // 返回结果
}
