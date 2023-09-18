package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // key: GET/POST, value: 单独的前缀树
	handlers map[string]HandlerFunc // 路径和路由方法
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 将 pattern 按照 "/" 进行拆分, 只允许一个 * 出现在路由匹配里
func parsePattern(pattern string) []string {
	rst := make([]string, 0)
	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		if part != "" {
			rst = append(rst, part)
		}
	}
	return rst
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// 树的根节点
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	// 返回了符合路径的叶节点
	parts := parsePattern(path)
	n := root.search(parts, 0) // 总是这一步搜索有问题

	if n == nil {
		return nil, nil
	}
	// 分析出路径参数
	params := make(map[string]string)
	oriParts := parsePattern(n.pattern)
	for idx, oriPart := range oriParts {
		var key string
		// 路由是 /a/:b/c, 请求的路径是 /a/b_value/c
		// 则 parmas={b: v_value}
		if strings.HasPrefix(oriPart, ":") {
			key = oriPart[1:]
			params[key] = parts[idx]
		}
		// 路由是 /a/*b, 请求的路径是 /a/bb/cc
		// params={b: "bb/cc"}
		if strings.HasPrefix(oriPart, "*") {
			key = oriPart[1:]
			params[key] = strings.Join(parts[idx:], "/")
			break
		}
	}
	return n, params
}

// 获取所有注册过的路由路径
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// 找出 context 对应的 handler 并返回 && 解析路径参数并存储在 context 中
func (r *router) handler(c *Context) {
	// 前缀树只是为了解析动态路由里的参数 && 判断请求路径是否被注册
	// 路径和路由方法还是由 map 保存
	n, params := r.getRoute(c.Method, c.Path)

	if n == nil {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	} else {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	}
}
