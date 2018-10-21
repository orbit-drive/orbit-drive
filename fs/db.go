package fs

import (
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

type FileStore map[string]string

var (
	Db *leveldb.DB
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

func GetFileStore() (FileStore, error) {
	store := make(FileStore)
	iter := Db.NewIterator(nil, nil)
	for iter.Next() {
		k := ToStr(iter.Key())
		switch k {
		case ROOT_KEY, CONFIG_KEY:
		default:
			v := ToStr(iter.Value())
			store[k] = v
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return FileStore{}, err
	}
	return store, nil
}

func RunGarbageCollection(s FileStore) error {
	b := new(leveldb.Batch)
	for k, _ := range s {
		b.Delete(ToByte(k))
	}
	err := Db.Write(b, nil)
	if err != nil {
		return err
	}
	return nil
}

func BatchPut(s FileStore) error {
	b := new(leveldb.Batch)
	for k, v := range s {
		b.Put(ToByte(k), ToByte(v))
	}
	return Db.Write(b, nil)
}
