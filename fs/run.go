package fs

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gogo/protobuf/proto"
	"github.com/orbit-drive/orbit-drive/fs/db"
	"github.com/orbit-drive/orbit-drive/fs/ipfs"
	"github.com/orbit-drive/orbit-drive/fs/p2p"
	"github.com/orbit-drive/orbit-drive/fs/sys"
	"github.com/orbit-drive/orbit-drive/fs/vtree"
	log "github.com/sirupsen/logrus"
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

func initWatcher(c *Config, vt *vtree.VTree) (*Watcher, error) {
	w, err := NewWatcher(c.Root)
	if err != nil {
		return &Watcher{}, err
	}
	log.WithField("path", c.Root).Info("Watching folder")
	dirPaths := vt.AllDirPaths()
	w.BatchAdd(dirPaths)
	go w.Start(vt)
	return w, nil
}

// Run is the main entry point for orbit drive p2p sync.
func Run(c *Config) {
	sys.Notify("Starting file sync!")
	defer sys.Alert("Stopping file sync!")

	log.WithField("node-addr", c.NodeAddr).Info("Initializing ipfs shell...")
	ipfs.InitShell(c.NodeAddr)

	go func() {
		log.Info("Initializing p2p connection to bootstrap nodes...")
		if err := p2p.InitConn(); err != nil {
			sys.Fatal(err.Error())
		}
		log.Info("p2p network connections successfully established!")
	}()

	log.Info("Initializing vtree...")
	vt, err := initVTree(c)
	if err != nil {
		sys.Fatal(err.Error())
	}
	log.Info("Vtree successfully initialized!")

	log.Info("Initializing watcher...")
	watcher, err := initWatcher(c, vt)
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
