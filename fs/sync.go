package fs

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Sync struct {
	Done chan bool
	Path string
}

func NewSync(p string, addr string) *Sync {
	return &Sync{
		Done: make(chan bool),
		Path: p,
	}
}

func (s *Sync) Stop() {
	s.Done <- true
}

func (s *Sync) Start() {
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
				log.Println("Write", op.String())
			case fsnotify.Create:
				go Upload(op.Name)
			case fsnotify.Rename:
				log.Println("Rename", op.Name)
			case fsnotify.Remove:
				log.Println("Remove", op.String())
			case fsnotify.Write:
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
