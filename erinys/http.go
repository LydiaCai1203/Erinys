package erinys

import (
	"encoding/json"
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
	pc := engine.peerclient[peer]
	return pc, peer
}

// peers: ["127.0.0.1:8001", ...]
func (engine *HTTPEngine) RegisterPeer(peers ...string) {
	for _, peer := range peers {
		engine.peers.RegisterPeer(peer)
		engine.peerclient[peer] = &PeerClient{
			baseURL:  fmt.Sprintf("http://%s", peer),
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
	// 获取对应 group 实例
	groupkey := fmt.Sprintf("%s-%s", parts[1], req.Host)
	g, ok := groups[groupkey]
	if !ok {
		fmt.Fprintf(w, "500 internel error: %s\n", req.URL)
		return
	}
	// 获取 key 对应的 value 值
	keyname := parts[2]
	v, _ := g.Get(keyname)
	// 返回 json 格式的数据
	resp := map[string]interface{}{keyname: v}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		http.Error(w, err.Error(), 500)
	}
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
