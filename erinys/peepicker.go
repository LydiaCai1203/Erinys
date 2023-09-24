package erinys

type PeerPicker interface {
	PickPeer(string) (*PeerClient, string)
}
