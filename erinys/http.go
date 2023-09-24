package erinys

import (
	"fmt"
	"net/http"
	"strings"
)

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
	parts := strings.Split(path, "/")
	if !isget || len(parts) == 3 {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}

	// 调用本地缓存，没有则调用源站数据
	// groupname := parts[1]
	// keyname := parts[2]
}

func (engine *HTTPEngine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
