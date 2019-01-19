package fs

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogo/protobuf/proto"
	"github.com/orbit-drive/orbit-drive/fs/api"
	"github.com/orbit-drive/orbit-drive/fs/db"
	"github.com/orbit-drive/orbit-drive/fs/sys"
	"github.com/orbit-drive/orbit-drive/fs/vtree"
)

func loadAndInitVTree(root string) (*vtree.VTree, error) {
	sources, err := db.GetSources()
	if err != nil {
		return &vtree.VTree{}, err
	}

	vt, err := vtree.NewVTree(root, sources)
	if err != nil {
		return &vtree.VTree{}, err
	}
	sources.Dump()
	return vt, nil
}

// Run is the main entry point for orbit drive sync mode by:
// (i) generating a virtual tree representation of the syncing folder.
// (ii) starts the backend hub for device synchronization.
// (iii) starts the watcher for file changes in the syncing folder.
// (iv) relaying all local changes to backend hub.
func Run(c *Config) {
	sys.Notify("Starting file sync!")
	defer sys.Alert("Stopping file sync!")
	api.InitShell(c.NodeAddr)

	vt, err := loadAndInitVTree(c.Root)
	if err != nil {
		sys.Fatal(err.Error())
	}

	hub := NewHub(c.HubAddr, c.AuthToken)
	go hub.Dial()
	go hub.SyncTree(vt)
	defer hub.Stop()

	dirPaths := vt.AllDirPaths()
	watcher := NewWatcher(c.Root)
	watcher.BatchAdd(dirPaths)
	go watcher.Start(vt)
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
			if err = hub.PushMsg(parsedPb); err != nil {
				log.Println(err)
			}
		case <-close:
			return
		}
	}
}
