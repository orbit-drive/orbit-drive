package db

import (
	"os"

	"github.com/orbit-drive/orbit-drive/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	// Db represents a connection to leveldb
	Db *leveldb.DB
)

// InitDb initialize the global Db instance located at (HOME_PATH/.orbit-drive/datastore)
func InitDb() error {
	cp := utils.GetConfigDir() + "/datastore"

	if !utils.PathExists(cp) {
		os.Mkdir(cp, os.ModePerm)
	}

	var err error
	Db, err = leveldb.OpenFile(cp, nil)
	return err
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
