package fs

import (
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

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

func BatchPut(m map[string]string) error {
	b := new(leveldb.Batch)
	for k, v := range m {
		b.Put([]byte(k), []byte(v))
	}
	return Db.Write(b, nil)
}
