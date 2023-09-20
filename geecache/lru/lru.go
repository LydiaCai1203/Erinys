/*
1. get(key) 访问节点
2. moveHead(key) 将节点挪至链表队首
3. add(key) 添加至链表队首
4. remove(key) 删除队尾元素
*/
package lru

import "container/list"

// 返回值所占内存大小
type Value interface {
	Len() int
}

type Cache struct {
	maxBytes  int64                         // 允许使用的最大内存
	nbytes    int64                         // 当前已使用的内存
	ll        *list.List                    // 双向链表
	cache     map[string]*list.Element      // map, value 指向双向链表的节点
	OnEvicted func(key string, value Value) // 记录被移除时的回调函数
}

type entry struct {
	key   string
	value Value
}

func New(maxBytes int64, OnEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		OnEvicted: OnEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

// O(1); map 中找到节点; 放到链表头部;
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		// 取出 ele.Value 转换为 entry 类型的指针
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除队尾 和 map 里的数据
func (c *Cache) RemoveTail() {
	// 返回链表最后一个元素
	ele := c.ll.Back()
	if ele != nil {
		// 删除链表的最后一个元素
		c.ll.Remove(ele)
		// 删除 map 中对应的 key
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		// key 所占字节数 + kv.Value.Len() 返回 value 字节大小
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 被删除了回调一下
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	// ele 是指向双端链表节点的指针
	if ele, ok := c.cache[key]; ok {
		// 移到表头
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 可能新增的 value 比原来的大
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 插入表头
		ele := c.ll.PushFront(&entry{key, value})
		// 插入 map
		c.cache[key] = ele
		// 更新当前内存占用
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 如果当前内存占用超过最大内存限制，删除表尾元素
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveTail()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
