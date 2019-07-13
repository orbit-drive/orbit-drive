package sync

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gogo/protobuf/proto"
	"github.com/orbit-drive/orbit-drive/config"
	"github.com/orbit-drive/orbit-drive/db"
	"github.com/orbit-drive/orbit-drive/ipfs"
	"github.com/orbit-drive/orbit-drive/p2p"
	"github.com/orbit-drive/orbit-drive/sys"
	"github.com/orbit-drive/orbit-drive/vtree"
	"github.com/orbit-drive/orbit-drive/watcher"
	log "github.com/sirupsen/logrus"
)

func initVTree() (*vtree.VTree, error) {
	log.Info("Initializing vtree...")

	s, err := db.GetSources()
	if err != nil {
		return nil, err
	}

	vt := vtree.NewVTree()
	if err := vt.Build(s); err != nil {
		return nil, err
	}
	s.Dump()
	log.Info("VTree successfully initialized!")
	return vt, nil
}

func initWatcher(vt *vtree.VTree) (*watcher.Watcher, error) {
	log.Info("Initializing watcher...")

	w, err := watcher.NewWatcher(config.GetRoot())
	if err != nil {
		return nil, err
	}
	log.WithField("path", config.GetRoot()).Info("Watching root folder to sync")
	dirPaths := vt.AllDirPaths()
	w.BatchAdd(dirPaths)
	go w.Start(vt)

	log.Info("Watcher initialized!")
	return w, nil
}

func initP2P(c *config.Config) {
	log.Info("Initializing p2p connection to bootstrap nodes...")
	if err := p2p.InitConn(c.P2PPort, c.SecretPhrase); err != nil {
		sys.Fatal(err.Error())
	}
	log.Info("p2p network connections successfully established!")
}

// Run is the main entry point for orbit drive p2p sync.
func Run() {
	sys.Notify("Starting file sync!")
	defer sys.Alert("Stopping file sync!")

	log.WithField("node-addr", config.GetNodeAddr()).Info("Initializing ipfs shell...")
	ipfs.InitShell(config.GetNodeAddr())

	// TODO: remove -> moving to web hub ???
	// go initP2P()

	vt, err := initVTree()
	if err != nil {
		sys.Fatal(err.Error())
	}
	log.WithField("hash", vt.MerkleHash()).Info("VTree loaded merkle hash")

	watcher, err := initWatcher(vt)
	if err != nil {
		sys.Fatal(err.Error())
	}
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
			p2p.GetMerkleHash()
		case <-close:
			return
		}
	}
}
