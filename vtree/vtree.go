package vtree

import (
	"path/filepath"
	"sync"

	"github.com/orbit-drive/orbit-drive/db"
	"github.com/orbit-drive/orbit-drive/pb"
	"github.com/orbit-drive/orbit-drive/utils"
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
func NewVTree(path string) *VTree {
	return &VTree{
		Head: &VNode{
			ID:     utils.ToByte(ROOTKEY),
			Path:   path,
			Type:   DirCode,
			Links:  []*VNode{},
			Source: &db.Source{},
		},
		state: make(chan State),
	}
}

// StateChanges returns the state channel of the VTree.
func (vt *VTree) StateChanges() <-chan State {
	return vt.state
}

// PushToState generates and sends a State struct to the state channel.
func (vt *VTree) PushToState(p string, op opCode) {
	vt.state <- State{Path: p, Op: op}
}

// Build is a wrapper around PopulateNodes to set flag or
// auto upload unsync files to ipfs network.
func (vt *VTree) Build(s db.Sources) error {
	return vt.PopulateNodes(s, true)
}

// PopulateNodes recursively populates the file tree structure
// starting from the head.
func (vt *VTree) PopulateNodes(s db.Sources, upload bool) error {
	return vt.Head.PopulateNodes(s, upload)
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
	isDir, err := utils.IsDir(path)
	if err != nil {
		return err
	}
	if isDir {
		n.SetAsDir()
		// Read file content and upload
		n.PopulateNodes(db.Sources{}, true)
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

func (vt *VTree) RootPath() string {
	return vt.Head.GetPath()
}

// AllDirPaths returns all the dir path in the vtree.
func (vt *VTree) AllDirPaths() []string {
	return vt.Head.AllDirPaths()
}

// MerkleHash returns the merkle root hash.
func (vt *VTree) MerkleHash() string {
	return vt.Head.MerkleHash()
}
