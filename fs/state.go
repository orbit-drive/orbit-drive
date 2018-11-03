package fs

import (
	"path/filepath"

	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/fs/db"
)

// NewFile traverse VTree to locate path parent dir and
// add a new vnode.
func NewFile(path string) error {
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
		n.Save()
	}
	return nil
}

// DeleteFile -> UnlinkChild -> remove from db
func DeleteFile(path string) error {
	return nil
}
