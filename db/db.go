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

	_, err := os.Stat(cp)
	if os.IsNotExist(err) {
		os.Mkdir(cp, os.ModePerm)
	}

	Db, err = leveldb.OpenFile(cp, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func Put(k []byte, v []byte) error {
	return Db.Put(k, v, nil)
}

func CloseDb() {
	Db.Close()
}
