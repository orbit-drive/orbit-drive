package db

import (
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/orbit-drive/orbit-drive/common"
)

var (
	// Db represents a connection to leveldb
	Db *leveldb.DB
)

// InitDb initialize the global Db instance located at (HOME_PATH/.orbit-drive/datastore)
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

// Put is a wrapper to leveldb Put func
func Put(k []byte, v []byte) error {
	return Db.Put(k, v, nil)
}

// Get is a wrapper to leveldb Get func
func Get(k []byte) ([]byte, error) {
	return Db.Get(k, nil)
}

// CloseDb is a wrapper to leveldb Close func
func CloseDb() {
	Db.Close()
}
