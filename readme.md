# Erinys

Golang 练手项目，一个简陋的仅仅支持字符串的分布式缓存服务器。

## 0. 目录结构
```markdown
.
├── consistenhash                
│   ├── consistenhash.go          // 一致性 hash 算法实现
│   └── consistenhash_test.go     // 一致性 hash 测试用例
├── erinys 
│   ├── erinys_test.go
│   ├── group.go                  // 组概念，每个缓存服务都可以有多个组，类似 redis 的 db
│   ├── http.go                   // HTTP 服务端支持用户访问获取数据
│   ├── peepicker.go              // 接口文件
│   └── peerclient.go             // HTTP 客户端，提供节点间的访问
├── go.mod
├── lru
│   ├── lru.go                    // LRU 算法实现
│   ├── lru_test.go               // LRU 测试用例
│   └── safelru.go                // 并发缓存支持
├── main.go
└── readme.md
```

## 1. 基本流程
```markdown
1. 用户请求
http://host1:port1/cache/<groupname>/<keyname>

2. 查询数据所在节点信息
当前被请求的缓存服务器 A 根据 keyname 寻找真实节点信息;
如果 keyname 所在节点是当前节点，则走步骤 3;
否则走步骤4;

3. 从本地获取数据
若本地缓存有 keyname 数据，直接获取返回;
否则，从源站请求数据(db 查询之类) 并 更新本地缓存;

4. 从远程获取数据
http://host2:port2/cache/<groupname>/<keyname>
```

## 2. QuickStart
```markdown
go run main.go
```

## 3. 待优化点
```markdown
1. 考虑支持多种数据格式
参考 redis...
2. 考虑支持机制防止缓存击穿
记录下哪些请求是正在运行的，当有重复请求进入，则等待未完成请求结束，一并返回结果
3. 考虑使用 protobuf 优化节点间的二进制通信
```
