package fs

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	Done chan bool
	Path string
}

func NewWatcher(p string) *Watcher {
	return &Watcher{
		Done: make(chan bool),
		Path: p,
	}
}

func (s *Watcher) Stop() {
	s.Done <- true
}

func (s *Watcher) Start() {
	ops, errs := startWatcher(s.Path)
	defer close(ops)
	defer close(errs)

	for {
		select {
		case op := <-ops:
			if op.Op.String() == "" {
				continue
			}
			switch op.Op {
			case fsnotify.Chmod:
				// Change in file permission
				log.Println("Write", op.String())
			case fsnotify.Create:
				// File created
				log.Println("Create", op.String())
			case fsnotify.Rename:
				// File renamed -> also called after create, write
				log.Println("Rename", op.String())
			case fsnotify.Remove:
				// File removed
				log.Println("Remove", op.String())
			case fsnotify.Write:
				// File modified or moved
				log.Println("Write", op.String())
			default:
				continue
			}
		case err := <-errs:
			log.Println(err)
		case <-s.Done:
			log.Println("Ipfsync stopped.")
			return
		}
	}
}

func startWatcher(p string) (chan fsnotify.Event, chan error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}

	ops := make(chan fsnotify.Event)
	errs := make(chan error)

	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				ops <- event
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				errs <- err
			}
		}
	}()

	if err = w.Add(p); err != nil {
		log.Fatalln(err)
	}

	return ops, errs
}
