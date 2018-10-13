package fs

import (
	"log"
	"os"
	"os/user"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	ConfigDb *leveldb.DB
)

func GetCurrentUsr() *user.User {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return usr
}

func GenConfigPath() string {
	usr := GetCurrentUsr()
	cp := usr.HomeDir + "/.ip-drive"

	_, err := os.Stat(cp)
	if os.IsNotExist(err) {
		os.Mkdir(cp, os.ModeDir)
	}

	return cp
}

func InitConfig() error {
	cp := GenConfigPath()

	var err error
	ConfigDb, err = leveldb.OpenFile(cp, nil)
	if err != nil {
		return err
	}

	return nil
}
