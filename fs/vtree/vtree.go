package vtree

import (
	"path/filepath"

	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/fs/db"
)

var (
	// VTree is the root pointer to the virtual tree of the file
	// structure being synchronized.
	VTree *VNode
)

// InitVTree initialize a new virtual tree (VTree) given an absolute path.
func InitVTree(path string, s db.Sources) error {
	VTree = &VNode{
		Path:   path, // To optimize here -> start with "/" not abs path
		Id:     common.ToByte(common.ROOT_KEY),
		Type:   DirCode,
		Links:  []*VNode{},
		Source: &db.Source{},
	}

	err := VTree.PopulateNodes(s)
	if err != nil {
		return err
	}
	return nil
}

func Find(path string) (*VNode, error) {
	return VTree.FindChildAt(path)
}

// NewFile traverse VTree to locate path parent dir and add a new vnode.
func Add(path string) error {
	dir := filepath.Dir(path)
	vn, err := VTree.FindChildAt(dir)
	if err != nil {
		return err
	}
	n := vn.NewVNode(path)
	isDir, err := common.IsDir(path)
	if err != nil {
		return err
	}
	if isDir {
		n.SetAsDir()
		// Read file content and upload
		n.PopulateNodes(db.Sources{})
	} else {
		n.SetAsFile()
		n.SaveSource()
	}
	return nil
}

// DeleteFile -> UnlinkChild -> remove from db
func Remove(path string) error {
	return nil
}
