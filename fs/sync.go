package fs

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gogo/protobuf/proto"
	"github.com/orbit-drive/orbit-drive/fs/config"
	"github.com/orbit-drive/orbit-drive/fs/db"
	"github.com/orbit-drive/orbit-drive/fs/ipfs"
	"github.com/orbit-drive/orbit-drive/fs/p2p"
	"github.com/orbit-drive/orbit-drive/fs/sys"
	"github.com/orbit-drive/orbit-drive/fs/vtree"
	log "github.com/sirupsen/logrus"
)

func initVTree() (*vtree.VTree, error) {
	sources, err := db.GetSources()
	if err != nil {
		return &vtree.VTree{}, err
	}

	vt, err := vtree.NewVTree(config.GetRootPath(), sources)
	if err != nil {
		return &vtree.VTree{}, err
	}
	sources.Dump()
	return vt, nil
}

func initWatcher(vt *vtree.VTree) (*Watcher, error) {
	w, err := NewWatcher(config.GetRootPath())
	if err != nil {
		return &Watcher{}, err
	}
	log.WithField("path", config.GetRootPath()).Info("Watching folder")
	dirPaths := vt.AllDirPaths()
	w.BatchAdd(dirPaths)
	go w.Start(vt)
	return w, nil
}

// Sync is the main entry point for orbit drive p2p sync.
func Sync() {
	sys.Notify("Starting file sync!")
	defer sys.Alert("Stopping file sync!")

	log.WithField("node-addr", config.GetNodeAddr()).Info("Initializing ipfs shell...")
	ipfs.InitShell(config.GetNodeAddr())

	go func() {
		log.Info("Initializing p2p connection to bootstrap nodes...")
		if err := p2p.InitConn(config.GetNID()); err != nil {
			sys.Fatal(err.Error())
		}
		log.Info("p2p network connections successfully established!")
	}()

	log.Info("Initializing vtree...")
	vt, err := initVTree()
	if err != nil {
		sys.Fatal(err.Error())
	}
	log.Info("Vtree successfully initialized!")

	log.Info("Initializing watcher...")
	watcher, err := initWatcher(vt)
	if err != nil {
		sys.Fatal(err.Error())
	}
	log.Info("Watcher initialized!")
	defer watcher.Stop()

	close := make(chan os.Signal, 2)
	signal.Notify(close, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case state := <-vt.StateChanges():
			log.WithFields(log.Fields{
				"path":      state.Path,
				"operation": state.Op,
			}).Info("vtree state change detected!")
			vtPb := vt.ToProto()
			parsedPb, err := proto.Marshal(vtPb)
			if err != nil {
				sys.Alert(err.Error())
			}
			log.WithField("byte-data", parsedPb).Info("vtree successfully parsed to pb!")
		case <-close:
			return
		}
	}
}
