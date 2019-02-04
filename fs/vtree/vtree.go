package vtree

import (
	"errors"
	"path/filepath"
	"sync"

	"github.com/orbit-drive/orbit-drive/fs/db"
	"github.com/orbit-drive/orbit-drive/fs/pb"
	"github.com/orbit-drive/orbit-drive/fsutil"
)

type opCode int64

const (
	ROOTKEY = "ROOT_TREE"

	// AddedOp represents the create operation
	AddedOp = iota
	// ModifiedOp represents the create operation
	ModifiedOp = iota
	// RemovedOp represents the remove operation
	RemovedOp = iota
)

// State represents a vtree state change.
type State struct {
	Path string
	Op   opCode
}

// VTree represents the file tree structure
type VTree struct {
	sync.Mutex

	// Head is the root pointer to the virtual tree of the file structure being synchronized.
	Head *VNode

	// State channel
	state chan State
}

// NewVTree initialize a new virtual tree (VTree) given an absolute path.
func NewVTree(rootPath string, s db.Sources) (*VTree, error) {
	isDir, err := fsutil.IsDir(rootPath)
	if err == nil && !isDir {
		err = errors.New("root path must be a directory")
	}
	if err != nil {
		return &VTree{}, err
	}
	vt := &VTree{
		Head: &VNode{
			Name:   filepath.Base(rootPath), // name of dir being sync.
			ID:     fsutil.ToByte(ROOTKEY),  // special id for root node.
			Type:   DirCode,
			Links:  []*VNode{}, // children
			parent: nil,
		},
		state: make(chan State),
	}
	if err = vt.PopulateNodes(s); err != nil {
		return &VTree{}, nil
	}
	return vt, nil
}

// StateChanges returns the state channel of the VTree.
func (vt *VTree) StateChanges() <-chan State {
	return vt.state
}

// PushToState generates and sends a State struct to the state channel.
func (vt *VTree) PushToState(p string, op opCode) {
	vt.state <- State{Path: p, Op: op}
}

// PopulateNodes recursively populates the file tree structure
// starting from the head.
func (vt *VTree) PopulateNodes(s db.Sources) error {
	return vt.Head.PopulateNodes(s)
}

// Find recursively traverse down the tree structure from the
// root head and returns the vnode corresponding the path.
func (vt *VTree) Find(path string) (*VNode, error) {
	return vt.Head.FindChildAt(path)
}

// Add traverse VTree to locate path parent dir and add a new vnode.
func (vt *VTree) Add(path string) error {
	vt.Lock()
	defer vt.Unlock()

	dir := filepath.Dir(path)
	vn, err := vt.Find(dir)
	if err != nil {
		return err
	}
	n := vn.NewVNode(path)
	isDir, err := fsutil.IsDir(path)
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
	vt.PushToState(path, AddedOp)
	return nil
}

// Remove -> UnlinkChild -> remove from db
func (vt *VTree) Remove(path string) error {
	vt.PushToState(path, RemovedOp)
	return nil
}

// ToProto parse a vtree to protobuf.
func (vt *VTree) ToProto() *pb.FSTree {
	return &pb.FSTree{
		Head: vt.Head.ToProto(),
	}
}

// AllDirPaths returns all the dir path in the vtree.
func (vt *VTree) AllDirPaths() []string {
	return vt.Head.AllDirPaths()
}
