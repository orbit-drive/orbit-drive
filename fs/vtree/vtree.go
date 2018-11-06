package vtree

import (
	"path/filepath"
	"sync"

	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/fs/db"
)

const (
	CreateOp = iota
	WriteOp  = iota
	RemoveOp = iota
)

// State represents a vtree state change.
type State struct {
	Path string
	Op   int64
}

type VTree struct {
	sync.Mutex
	// VTree is the root pointer to the virtual tree of the file
	// structure being synchronized.
	Head *VNode

	// State channel
	state chan State
}

// InitVTree initialize a new virtual tree (VTree) given an absolute path.
func NewVTree(path string, s db.Sources) (*VTree, error) {
	vt := &VTree{
		Head: &VNode{
			Path:   path,
			Id:     common.ToByte(common.ROOT_KEY),
			Type:   DirCode,
			Links:  []*VNode{},
			Source: &db.Source{},
		},
		state: make(chan State),
	}

	err := vt.PopulateNodes(s)
	if err != nil {
		return &VTree{}, nil
	}
	return vt, nil
}

func (vt *VTree) StateChanges() <-chan State {
	return vt.state
}

func (vt *VTree) PushToState(p string, op int64) {
	vt.state <- State{Path: p, Op: op}
}

func (vt *VTree) PopulateNodes(s db.Sources) error {
	return vt.Head.PopulateNodes(s)
}

func (vt *VTree) Find(path string) (*VNode, error) {
	return vt.Head.FindChildAt(path)
}

// Add traverse VTree to locate path parent dir and add a new vnode.
func (vt *VTree) Add(path string) error {
	vt.Lock()
	defer vt.Unlock()

	dir := filepath.Dir(path)
	vn, err := vt.Head.FindChildAt(dir)
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
func (vt *VTree) Remove(path string) error {
	return nil
}
