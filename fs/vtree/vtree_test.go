package vtree

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/orbit-drive/orbit-drive/fs/db"
)

const (
	TESTDATA_DIRNAME = "testdata"

	TESTDATA_ROOTHASH = "111664f60b22a9baef139fa783accae108e0a3ddd1403437c6d280b37f8ffbd7"
)

func setupTestVTree() (*VTree, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	testDataPath := filepath.Join(path, TESTDATA_DIRNAME)
	return NewVTree(testDataPath), nil
}

func TestVTreeInit(t *testing.T) {
	vt, err := setupTestVTree()
	if err != nil {
		t.Error(err)
	}

	vt.PopulateNodes(make(db.Sources), false)

	head := vt.Head
	if len(head.Links) != 2 {
		t.Errorf("Expected %d vnodes, got: %d", 2, len(head.Links))
	}

	head.SortLinksByID()
	folder1 := head.Links[0]
	file1 := head.Links[1]

	if !folder1.IsDir() {
		t.Error("folder1 should be a dir not a file.")
	}
	if file1.IsDir() {
		t.Error("file1 should be a file not a dir.")
	}

	if folder1.LinksCount() != 1 {
		t.Errorf("folder1 should have %d vnode, got: %d", 1, folder1.LinksCount())
	}
	folder1Child := folder1.Links[0]
	if folder1Child.GetName() != "file2" {
		t.Errorf("Expected folder1 child filename to be %s, got: %s", "file2", folder1Child.GetName())
	}
}

func TestMerkleHash(t *testing.T) {
	vt, err := setupTestVTree()
	if err != nil {
		t.Error(err)
	}

	var rootHash string
	rootHash = vt.MerkleHash()
	if rootHash != "" {
		// Should be empty since vtree is not populated yet.
		t.Errorf("Expected empty root hash got: %s", rootHash)
	}

	// Populate vtree nodes.
	vt.PopulateNodes(make(db.Sources), false)
	rootHash = vt.MerkleHash()
	if rootHash != TESTDATA_ROOTHASH {
		t.Errorf("Expected empty root hash got: %s", rootHash)
	}

}
