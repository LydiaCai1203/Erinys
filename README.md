# myRPC

## 1. Gee

### 1.1. 动态路由

```markdown
借助 前缀树 实现，前缀树有已注册路由信息组成，主要用于判断请求路径是否存在 && 路径参数解析;
借助 map 存储注册路径和路由方法的映射关系，将请求路径的参数部分替换再去 map 中查找对应的函数;
```

### 1.2. 路由分组

```markdown
1. 例子
/post: 该前缀开头的路由匿名可访问
/admin: 该前缀开头的路由需要鉴权
/api: 该前缀开头的路由是 RESTful 接口，可以对接第三方平台，需要三方平台鉴权

2. 中间件
作用在 /post 分组上的中间件也会作用在其子分组上，子分组也可以单独应用自己的中间件
```

### 1.3. 中间件

```markdown
// Next 函数
func (c *Context) Next() {
    c.index++
    s := len(c.handlers)
    for ; c.index < s; c.index++ {
        c.handlers[c.index](c)
    }
}

// 有 A、B 两个中间件函数
func A(c *Context) {
    part1      // 执行路由函数前调用
    c.Next()
    part2      // 执行路由函数后调用
}

func B(c *Context) {
    part3      // 执行路由函数前调用
    c.Next()
    part4      // 执行路由函数后调用
}

// C 是路由函数
// 使用 Next 调用 handlers 里的函数
// 顺序: part1 -> part3 -> C -> part4 -> part2
handlers := []HandleFunc{A, B, C}
```

## 2. GeeCache

```markdown
1. 概念
groupcache，缓存数据的分布式缓存库，最初由 Google 设计和开发。
GeeCache 基本模仿了 groupcache 的实现。

2. 支持特性
单机缓存和基于 HTTP 的分布式缓存;
LRU 缓存策略;
使用 Go 锁机制避免缓存击穿;
使用一致性哈希选择节点，实现负载均衡;
使用 protobuf 优化节点间二进制通信;
```

### 2.1. LRU 缓存淘汰策略

**1. FIFO**
```markdown
先进先出。最早添加的记录，不在被使用的可能性最大。
实现时创建一个队列，队尾添加，内存不够时淘汰队首。

+ 缺点
最早记录但却最常访问的数据，会被频繁地添加和淘汰，导致命中率低。
```

**2. LFU**
```markdown
淘汰掉缓存中访问频率最低的记录。
实现时维护一个按照访问次数排序的队列，每次访问，访问次数加 1，队列重新排序。淘汰时选择访问次数最少的淘汰即可。

+ 缺点
维护每个记录的访问次数，对内存的消耗较高;
如果数据的访问模式发生变化，LFU 需要较长时间适应;
```

**3. LRU**
```markdown
淘汰掉最近最少使用的记录;
实现时维护一个队列，如果某条记录被访问了，则移动到队尾，淘汰队首数据即可。
```

### 2.2. 单机并发缓存

```golang
type Value interface {
    Len() int
}

type entry struct {
    key   string
    value Value
}

type ByteView struct {
    b []byte                                // 存储真实的缓存值
}

type Cache struct {
    maxBytes  int64                         // 允许使用的最大内存
    nbytes    int64                         // 当前已使用的内存
    ll        *list.List                    // 双向链表
    cache     map[string]*list.Element      // map, value 指向双向链表的节点
    OnEvicted func(key string, value Value) // 记录被移除时的回调函数
}

type cache struct {
    mu         sync.Mutex
    lru        *lru.Cache
    cacheBytes int64
}

map 里存的 value 是指向双向链表的节点的指针, *list.Element。
双向链表的节点是 entry 类型, 存储了 k/v 值, v 值类型是 Value 接口类型的。
Value 接口类型可以转化尾 ByteView 类型。
cache 类型是一个加锁的 lru 模型，适合单机并发存储。
```

### 2.3. 一致性哈希

```markdown
1. 算法描述
一致性哈希算法就是将 key 映射到 2^32 的空间中，然后将这些数字首尾相连，形成一个环。
计算 节点/机器 的哈希值(通常是节点的名称、IP 地址、编号)，放置在环上。
计算 key 的哈希值，放置在环上，然后顺时针找到第一个节点，就是应选取的 节点/机器。

2. 数据倾斜问题
如果服务器节点过少容易引起 key 的倾斜(key 的分布不均)。
为了解决此类问题引入 虚拟节点。一个真实的节点对应多个虚拟节点。
比如一个真实节点，对应三个虚拟节点，计算虚拟节点的哈希值并放置在环上。
计算 key 的哈希值，在环上顺时针寻找到应选取的虚拟节点，然后找到虚拟节点对应的真实节点，存入即可。
```
