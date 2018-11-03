package db

import (
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/wlwanpan/orbit-drive/common"
)

var (
	Db *leveldb.DB
)

func InitDb() {
	cp := common.GetHomeDir() + "/.orbit-drive/datastore"

	if !common.PathExists(cp) {
		os.Mkdir(cp, os.ModePerm)
	}

	var err error
	Db, err = leveldb.OpenFile(cp, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func Put(k []byte, v []byte) error {
	return Db.Put(k, v, nil)
}

func Get(k []byte) ([]byte, error) {
	return Db.Get(k, nil)
}

func CloseDb() {
	Db.Close()
}
