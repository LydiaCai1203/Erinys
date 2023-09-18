/*
gorouter 支持在路由规则中嵌入正则表达式;
httprouter 不支持正则表达，gin 曾经时候，后又放弃;
Trie 树(前缀树) 是实现动态路由最常见的数据结构，每一个节点的所有子节点都拥有相同的前缀，这种结构适合实现动态路由匹配;

动态路由匹配规则:
1. 参数匹配 ":"
/p/:lang/doc 可以匹配 /p/c/doc 和 /p/go/doc
2. 通配 "*"
/static/*filepath 可以匹配 /static/fav.ico 和 /static/js/jQuery.js
*/

package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由, /p/:lang, 只有叶子节点有
	part     string  // 路由中的一部分, :lang
	children []*node // 子节点
	isWild   bool    // 是否属于模糊匹配节点
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

func (n *node) matchChild(part string) *node {
	// 只能找当前节点 n 的孩子节点里是否有符合的，没有就认为没有
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 因为有 * 通配符的缘故，所以可能找到多个孩子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 这样写插入 /a/*b/c, 会出错; 插入 /a/*b 不会出错;
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*"),
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 查找 parts 所指路径是否存在于树中
func (n *node) search(parts []string, height int) *node {
	// 找到了叶子节点 或 匹配到了通配符节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 存了 /a/b/c 但是搜的是 /a/b, n.pattern 就是 ""
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

// 遍历出所有的路径
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
