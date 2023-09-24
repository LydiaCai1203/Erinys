package erinys

import (
	"erinys/consistenhash"
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
	basepath   string
	router     map[string]HandlerFunc
	peers      *consistenhash.PeerPool // PeerPool 实例
	peerclient map[string]*PeerClient  // 真实节点名称 与 PeerClient 的映射
}

func NewHTTPEngine(
	basepath string,
	replicas int,
	fn consistenhash.HashFunc,
) *HTTPEngine {
	return &HTTPEngine{
		basepath:   basepath,
		router:     make(map[string]HandlerFunc),
		peers:      consistenhash.NewPeerPool(replicas, fn),
		peerclient: make(map[string]*PeerClient),
	}
}

func (engine *HTTPEngine) PickPeer(key string) (*PeerClient, string) {
	peer := engine.peers.GetPeerByKey(key)
	pc := engine.peerclient[key]
	return pc, peer
}

// peers: ["127.0.0.1:8001", ...]
func (engine *HTTPEngine) RegisterPeer(peers ...string) {
	for _, peer := range peers {
		engine.peerclient[peer] = &PeerClient{
			baseURL:  fmt.Sprintf("http:%s", peer),
			basepath: "cache",
		}
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
		fmt.Fprintf(w, "500 internel error: %s\n", req.URL)
		return
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
