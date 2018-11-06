package fs

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wlwanpan/orbit-drive/fs/api"
	"github.com/wlwanpan/orbit-drive/fs/db"
	"github.com/wlwanpan/orbit-drive/fs/sys"
	"github.com/wlwanpan/orbit-drive/fs/vtree"
)

func Run(c *Config) {
	sys.Notify("Starting file sync!")
	defer sys.Alert("Stopping file sync!")
	api.InitShell(c.NodeAddr)

	sources, err := db.GetSources()
	if err != nil {
		sys.Fatal(err.Error())
	}

	vt, err := vtree.NewVTree(c.Root, sources)
	if err != nil {
		sources.Dump()
	}

	hub, err := NewHub("localhost:4000") // To move to env
	if err != nil {
		sys.Alert(err.Error())
	} else {
		go hub.Sync(vt)
		defer hub.Stop()
	}

	watcher := NewWatcher(c.Root)
	go watcher.Start(vt)
	defer watcher.Stop()

	close := make(chan os.Signal, 2)
	signal.Notify(close, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case state := <-vt.StateChanges():
			log.Println(state.Path)
		case <-close:
			return
		}
	}
}
