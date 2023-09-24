package consistenhash

import (
	"fmt"
	"testing"
)

func TestConsistenHash(t *testing.T) {
	pool := NewPeerPool(5, nil)
	peers := []string{"A", "B", "C", "D"}
	pool.RegisterPeer(peers...)
	peer := pool.GetPeerByKey("name")
	fmt.Println(peer)
}
