# myRPC

## 1. 动态路由
```markdown
借助 前缀树 实现，前缀树有已注册路由信息组成，主要用于判断请求路径是否存在 && 路径参数解析;
借助 map 存储注册路径和路由方法的映射关系，将请求路径的参数部分替换再去 map 中查找对应的函数;
```

## 2. 路由分组
```markdown
1. 例子
/post: 该前缀开头的路由匿名可访问
/admin: 该前缀开头的路由需要鉴权
/api: 该前缀开头的路由是 RESTful 接口，可以对接第三方平台，需要三方平台鉴权

2. 中间件
作用在 /post 分组上的中间件也会作用在其子分组上，子分组也可以单独应用自己的中间件
```

## 3. 中间件
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
