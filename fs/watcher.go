package fs

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"sync"

	"github.com/fsnotify/fsnotify"
)

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
	populateWatchlist(n, w.Path)
	w.Notifier = n
}

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
func populateWatchlist(w *fsnotify.Watcher, p string) {
	var wg sync.WaitGroup
	if err := w.Add(p); err != nil {
		w.Errors <- err
	}
	fmt.Println("Watching: ", p)

	files, err := ioutil.ReadDir(p)
	if err != nil {
		w.Errors <- err
		return
	}

	for _, f := range files {
		if f.IsDir() {
			wg.Add(1)
			go func() {
				filepath := path.Join(p, f.Name())
				populateWatchlist(w, filepath)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}
