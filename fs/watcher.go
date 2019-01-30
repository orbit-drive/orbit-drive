package fs

import (
	"github.com/fsnotify/fsnotify"
	"github.com/orbit-drive/orbit-drive/fs/db"
	"github.com/orbit-drive/orbit-drive/fs/sys"
	"github.com/orbit-drive/orbit-drive/fs/vtree"
	"github.com/orbit-drive/orbit-drive/fsutil"
	log "github.com/sirupsen/logrus"
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
func NewWatcher(p string) (*Watcher, error) {
	n, err := fsnotify.NewWatcher()
	if err != nil {
		return &Watcher{}, err
	}
	w := &Watcher{
		Done:     make(chan bool),
		Path:     p,
		Notifier: n,
	}
	w.AddToWatchList(w.Path)
	return w, nil
}

// AddToWatchList adds path to watch.
func (w *Watcher) AddToWatchList(p string) {
	if err := w.Notifier.Add(p); err != nil {
		w.Notifier.Errors <- err
	}
}

// BatchAdd adds multiple paths to watcher.
func (w *Watcher) BatchAdd(paths []string) {
	for _, path := range paths {
		w.AddToWatchList(path)
		log.WithField("path", path).Info("New path added to watcher!")
	}
}

// RemoveFromWatchList traverse a path and remove all nested dir
// path from the notifier watch list.
func (w *Watcher) RemoveFromWatchList(p string) {
	if err := w.Notifier.Remove(p); err != nil {
		w.Notifier.Errors <- err
	}
}

// Start initialize watcher notifier and check for the notifier
// Event channel and Errors channel.
func (w *Watcher) Start(vt *vtree.VTree) {
	for {
		select {
		case e := <-w.Notifier.Events:
			if !validEvent(e) {
				continue
			}
			switch e.Op {
			case fsnotify.Create:
				createHandler(w, vt, e.Name)
			case fsnotify.Write:
				writeHandler(w, vt, e.Name)
			case fsnotify.Remove: // manually remove from folder triggers fsnotify.Rename
				removeHandler(w, vt, e.Name)
			default:
				log.WithField("event", e.String()).Warn("Watcher unhandled event")
				continue
			}
		case err := <-w.Notifier.Errors:
			sys.Alert(err.Error())
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

func createHandler(w *Watcher, vt *vtree.VTree, p string) {
	log.WithField("path", p).Info("Watcher detected file op: create")
	if err := vt.Add(p); err != nil {
		sys.Alert(err.Error())
		return
	}
	if isDir, _ := fsutil.IsDir(p); isDir {
		w.AddToWatchList(p) // TODO: Figure out how to get all new dir paths from vt.Add
	}
}

func writeHandler(w *Watcher, vt *vtree.VTree, p string) {
	log.WithField("path", p).Info("Watcher detected file op: write")
	vn, err := vt.Find(p)
	if err != nil {
		sys.Alert(err.Error())
		return
	}
	source := db.NewSource(p)
	if err := vn.UpdateSource(source); err != nil {
		log.Warn(err)
	}
	vt.PushToState(vn.Path, vtree.ModifiedOp)
}

func removeHandler(w *Watcher, vt *vtree.VTree, p string) {
	log.WithField("path", p).Info("Watcher detected file op: remove")
	vt.Remove(p)
	if isDir, _ := fsutil.IsDir(p); isDir {
		w.RemoveFromWatchList(p) // TODO: Figure out how to get all dir paths removed from vt.Remove
	}
}

func validEvent(e fsnotify.Event) bool {
	if e.Op.String() == "" {
		return false
	}
	return !fsutil.IsHidden(e.Name)
}
