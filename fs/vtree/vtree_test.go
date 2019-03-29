package vtree

import (
	"testing"
)

func TestMerkleHash(t *testing.T) {
	vt := &VTree{
		Head: &VNode{
			Path: "/test",
			Type: DirCode,
		},
	}

	rootHash := vt.MerkleHash()
	if rootHash != "" {
		t.Errorf("Expected empty root hash got: %s", rootHash)
	}
}
