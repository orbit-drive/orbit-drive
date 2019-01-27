package fs

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogo/protobuf/proto"
	"github.com/orbit-drive/orbit-drive/fs/db"
	"github.com/orbit-drive/orbit-drive/fs/ipfs"
	"github.com/orbit-drive/orbit-drive/fs/p2p"
	"github.com/orbit-drive/orbit-drive/fs/sys"
	"github.com/orbit-drive/orbit-drive/fs/vtree"
)

func initVTree(c *Config) (*vtree.VTree, error) {
	sources, err := db.GetSources()
	if err != nil {
		return &vtree.VTree{}, err
	}

	vt, err := vtree.NewVTree(c.Root, sources)
	if err != nil {
		return &vtree.VTree{}, err
	}
	sources.Dump()
	return vt, nil
}

func initHub(c *Config, vt *vtree.VTree) *Hub {
	h := NewHub(c.HubAddr, c.AuthToken)
	go h.Dial()
	go h.SyncTree(vt)
	return h
}

func initWatcher(c *Config, vt *vtree.VTree) *Watcher {
	w := NewWatcher(c.Root)
	dirPaths := vt.AllDirPaths()
	w.BatchAdd(dirPaths)
	go w.Start(vt)
	return w
}

// Run is the main entry point for orbit drive sync mode by:
// (i) generating a virtual tree representation of the syncing folder.
// (ii) starts the backend hub for device synchronization.
// (iii) starts the watcher for file changes in the syncing folder.
// (iv) relaying all local changes to backend hub.
func Run(c *Config) {
	sys.Notify("Starting file sync!")
	defer sys.Alert("Stopping file sync!")
	ipfs.InitShell(c.NodeAddr)

	vt, err := initVTree(c)
	if err != nil {
		sys.Fatal(err.Error())
	}

	// Moving hub to p2p connection to sync device.
	// hub := initHub(c, vt)
	// defer hub.Stop()

	go func() {
		if err := p2p.InitConn(); err != nil {
			sys.Fatal(err.Error())
		}
	}()

	watcher := initWatcher(c, vt)
	defer watcher.Stop()

	close := make(chan os.Signal, 2)
	signal.Notify(close, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case state := <-vt.StateChanges():
			log.Println(state.Path)
			vtPb := vt.ToProto()
			parsedPb, err := proto.Marshal(vtPb)
			if err != nil {
				sys.Alert(err.Error())
			}
			log.Println(parsedPb)
			// if err = hub.PushMsg(parsedPb); err != nil {
			// 	log.Println(err)
			// }
		case <-close:
			return
		}
	}
}
