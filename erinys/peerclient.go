package erinys

import (
	"erinys/lru"
	"fmt"
	"io"
	"net/http"
)

type Byte []byte

func (b Byte) Len() int64 {
	return int64(len(b))
}

type PeerClient struct {
	baseURL  string // http://127.0.0.1:8080
	basepath string // cache
}

func (pc *PeerClient) Get(group string, key string) (lru.Value, error) {
	// http://127.0.0.1:8080/cache/<groupname>/<key>
	url := fmt.Sprintf("%s%s/%s/%s", pc.baseURL, pc.basepath, group, key)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return Byte(bytes), err
}
