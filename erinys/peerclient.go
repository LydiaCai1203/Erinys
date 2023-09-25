package erinys

import (
	"encoding/json"
	"erinys/lru"
	"fmt"
	"io"
	"net/http"
)

type PeerClient struct {
	baseURL  string // http://127.0.0.1:8080
	basepath string // cache
}

func (pc *PeerClient) Get(group string, key string) (lru.Value, error) {
	// http://127.0.0.1:8080/cache/<groupname>/<key>
	url := fmt.Sprintf("%s/%s/%s/%s", pc.baseURL, pc.basepath, group, key)
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}
	bytes, _ := io.ReadAll(res.Body)
	m := make(map[string]string)
	err = json.Unmarshal([]byte(bytes), &m)
	if err != nil {
		return nil, nil
	}
	return String(m[key]), nil
}
