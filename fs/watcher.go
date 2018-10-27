package fs

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Callback func(string)

// Watcher is a wrapper to fsnotify watcher and represents
// a path to watch for usr changes.
type Watcher struct {
	// Done channel used to indicate when to stop the watcher
	Done chan bool

	// Path is absolute path to watch.
	Path string

	// Notifier holds the fs watcher
	Notifier *fsnotify.Watcher
}

// NewWatcher initialize a new Watcher.
func NewWatcher(p string) *Watcher {
	return &Watcher{
		Done: make(chan bool),
		Path: p,
	}
}

// InitNotifier initialize the wrapped fs watcher and populate its watch list.
func (w *Watcher) InitNotifier() {
	n, err := fsnotify.NewWatcher()
	if err != nil {
		// To figure out how to deal with error here <<
		fmt.Println(err)
	}
	w.Notifier = n
	w.PopulateWatchList(w.Path)
}

// PopulateWatchList traverse a path and add all nested dir
// path to the notifier watch list.
func (w *Watcher) PopulateWatchList(p string) {
	cb := func(path string) {
		if err := w.Notifier.Add(path); err != nil {
			w.Notifier.Errors <- err
		}
	}
	traversePath(p, cb)
}

// RemoveFromWatchList traverse a path and remove all nested dir
// path from the notifier watch list.
func (w *Watcher) RemoveFromWatchList(p string) {
	cb := func(path string) {
		if err := w.Notifier.Remove(path); err != nil {
			w.Notifier.Errors <- err
		}
	}
	traversePath(p, cb)
}

// Start initialize watcher notifier and check for the notifier
// Event channel and Errors channel.
func (w *Watcher) Start() {
	w.InitNotifier()

	for {
		select {
		case e := <-w.Notifier.Events:
			if e.Op.String() == "" {
				continue
			}
			switch e.Op {
			case fsnotify.Chmod:
				// Change in file permission
				log.Println("Write", e.String())
			case fsnotify.Create:
				// File created
				log.Println("Create", e.String())
				err := NewFile(e.Name)
				if err != nil {
					log.Println(err)
				}
			case fsnotify.Rename:
				// File renamed -> also called after create, write
				log.Println("Rename", e.String())
			case fsnotify.Remove:
				// File removed
				log.Println("Remove", e.String())
				w.RemoveFromWatchList(e.Name)
			case fsnotify.Write:
				// File modified or moved
				log.Println("Write", e.String())
			default:
				continue
			}
		case err := <-w.Notifier.Errors:
			log.Println(err)
		case <-w.Done:
			log.Println("Ipfsync stopped.")
			return
		}
	}
}

// Stop close the fs watcher and triggers the Done channel.
func (w *Watcher) Stop() {
	w.Notifier.Close()
	w.Done <- true
}

// populateWatchlist is a recursive func that traverse all nested
// dir of the given path and adds them to the fsnotify watch list.
func traversePath(p string, cb Callback) {
	var wg sync.WaitGroup

	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Println(err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			wg.Add(1)
			go func() {
				filepath := path.Join(p, f.Name())
				traversePath(filepath, cb)
				wg.Done()
			}()
		}
	}

	wg.Wait()
	cb(p)
}
