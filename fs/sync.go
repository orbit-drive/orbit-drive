package fs

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/syndtr/goleveldb/leveldb"
)

type IpfsSync struct {
	Done  chan bool
	Path  string
	Shell *shell.Shell
	Db    *leveldb.DB
}

func NewIpfsSync(p string, addr string) *IpfsSync {
	db, err := leveldb.OpenFile(p, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &IpfsSync{
		Done:  make(chan bool),
		Path:  p,
		Shell: shell.NewShell(addr),
		Db:    db,
	}
}

func (is *IpfsSync) Stop() {
	is.Done <- true
}

func (is *IpfsSync) Start() {
	ops, errs := startWatcher(is.Path)
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
				go is.Upload(op.Name)
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
		case <-is.Done:
			log.Println("Ipfsync stopped.")
			return
		}
	}
}

func (is *IpfsSync) Upload(p string) error {
	fi, err := os.Stat(p)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		log.Println("Uploading dir: ", p)
		cid, err := is.Shell.AddDir(p)
		if err != nil {
			return err
		}
		log.Println("Uploaded dir: ", cid)
		return nil
	}

	f, _ := os.Open(p)
	log.Println("Uploading file: ", p)
	cid, err := is.Shell.Add(f)
	if err != nil {
		return err
	}

	err = is.Db.Put([]byte(p), []byte(cid), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Uploaded file: ", cid)

	return nil
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
