package p2p

import (
	"errors"
)

const (
	MerkleHashRequest = "MerkleHashRequest"
)

var (
	ErrLNodeNotInitialized = errors.New("p2p: local node not initialized")
)

func sendRequest(method string) error {
	if lnode == nil {
		return ErrLNodeNotInitialized
	}
	lnode.Request(method)
	return nil
}

func GetMerkleHash() error {
	return sendRequest(MerkleHashRequest)
}
