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

	if err = vtree.InitVTree(c.Root, sources); err != nil {
		sources.Dump()
	}

	hub := NewHub()
	go hub.Start()
	defer hub.Stop()

	watcher := NewWatcher(c.Root)
	go watcher.Start()
	defer watcher.Stop()

	error := make(chan error)
	close := make(chan os.Signal, 2)
	signal.Notify(close, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case hs := <-hub.State:
			log.Println(hs.Path + hs.Op)
		case ws := <-watcher.State:
			log.Println(ws.Path + ws.Op)
			switch ws.Op {
			case "Create":
				error <- vtree.Add(ws.Path)
			case "Write":
				vn, err := vtree.Find(ws.Path)
				if err != nil {
					error <- err
					continue
				}
				source := db.NewSource(ws.Path)
				if vn.Source.IsSame(source) {
					continue
				}
				vn.Source = source
				vn.SaveSource()
			default:
			}
		case err := <-error:
			if err != nil {
				sys.Alert(err.Error())
			}
		case <-close:
			return
		}
	}
}
