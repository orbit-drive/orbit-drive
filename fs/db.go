package fs

import (
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	Db         *leveldb.DB
	SavedFiles map[string]string
)

func InitDb() {
	cp := GetHomeDir() + "/.orbit-drive/datastore"

	_, err := os.Stat(cp)
	if os.IsNotExist(err) {
		os.Mkdir(cp, os.ModePerm)
	}

	Db, err = leveldb.OpenFile(cp, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func InitSavedFiles() error {
	SavedFiles = map[string]string{}
	iter := Db.NewIterator(nil, nil)
	for iter.Next() {
		k := string(iter.Key()[:])
		switch k {
		case ROOT_KEY, CONFIG_KEY:
		default:
			v := string(iter.Value()[:])
			SavedFiles[k] = v
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return err
	}
	return nil
}

func RunGarbageCollect() error {
	b := new(leveldb.Batch)
	for k, _ := range SavedFiles {
		b.Delete([]byte(k))
	}
	err := Db.Write(b, nil)
	if err != nil {
		return err
	}
	SavedFiles = map[string]string{}
	return nil
}

func BatchPut(m map[string]string) error {
	b := new(leveldb.Batch)
	for k, v := range m {
		b.Put([]byte(k), []byte(v))
	}
	return Db.Write(b, nil)
}
