package erinys

import (
	"erinys/lru"
	"fmt"
	"net/http"
	"strings"
)

type String string

func (s String) Len() int64 {
	return int64(len(s))
}

type HandlerFunc func(w http.ResponseWriter, req *http.Request)

// 用户可以通过 HTTP 请求的方式来访问缓存服务器
type HTTPEngine struct {
	basepath string
	router   map[string]HandlerFunc
}

func NewHTTPEngine(basepath string) *HTTPEngine {
	return &HTTPEngine{
		basepath: basepath,
		router:   make(map[string]HandlerFunc),
	}
}

// basepath/<groupname>/<keyname>
func (engine *HTTPEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	isget := strings.HasPrefix(path, engine.basepath)
	parts := parsePath(path)
	if !isget || len(parts) != 3 {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
		return
	}

	groupname := parts[1]
	keyname := parts[2]
	g, ok := groups[groupname]
	if !ok {
		g = NewGroup(
			parts[1],
			GetterFunc(
				func(key string) (lru.Value, error) {
					m := map[string]string{
						"key1": "value1",
						"key2": "value2",
						"key3": "value3",
					}
					v, ok := m[key]
					if !ok {
						return nil, fmt.Errorf("%s not exit", key)
					}
					vv := String(v)
					return vv, nil
				}),
			2<<3,
		)
	}
	v, _ := g.Get(keyname)
	fmt.Fprintf(w, "%s-%v", keyname, v)
}

func (engine *HTTPEngine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func parsePath(path string) []string {
	rst := make([]string, 0)
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part == "" {
			continue
		}
		rst = append(rst, part)
	}
	return rst
}
