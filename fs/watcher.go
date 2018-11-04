package fs

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/fs/db"
	"github.com/wlwanpan/orbit-drive/fs/vtree"
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
	w.AddToWatchList(w.Path)
}

// AddToWatchList traverse a path and add all nested dir
// path to the notifier watch list.
func (w *Watcher) AddToWatchList(p string) {
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
		log.Println("Removed from watch list: ", path)
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
				chmodHandler(w, e.Name)
			case fsnotify.Create:
				createHandler(w, e.Name)
			case fsnotify.Rename:
				renameHandler(w, e.Name)
			case fsnotify.Remove:
				removeHandler(w, e.Name)
			case fsnotify.Write:
				writeHandler(w, e.Name)
			default:
				continue
			}
		case err := <-w.Notifier.Errors:
			log.Println(err)
		case <-w.Done:
			return
		}
	}
}

// Stop close the fs watcher and triggers the Done channel.
func (w *Watcher) Stop() {
	w.Notifier.Close()
	w.Done <- true
}

func chmodHandler(w *Watcher, p string) {
	log.Println("Write", p)
}

func createHandler(w *Watcher, p string) {
	log.Println("Create", p)
	err := vtree.Add(p)
	if err != nil {
		log.Println(err)
		return
	}
	if isDir, _ := common.IsDir(p); isDir {
		w.AddToWatchList(p)
	}
}

func renameHandler(w *Watcher, p string) {
	log.Println("Rename", p)
}

func removeHandler(w *Watcher, p string) {
	log.Println("Remove", p)
	if isDir, _ := common.IsDir(p); isDir {
		w.RemoveFromWatchList(p)
	}
}

func writeHandler(w *Watcher, p string) {
	log.Println("Write", p)
	vn, err := vtree.Find(p)
	if err != nil {
		log.Println(err)
		return
	}
	source := db.NewSource(p)
	if vn.Source.IsSame(source) {
		return
	}
	vn.Source = source
	vn.SaveSource()
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
		if !f.IsDir() {
			continue
		}
		wg.Add(1)
		go func() {
			filepath := path.Join(p, f.Name())
			traversePath(filepath, cb)
			wg.Done()
		}()
	}

	wg.Wait()
	cb(p)
}
