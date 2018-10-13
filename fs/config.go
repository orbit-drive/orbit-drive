package fs

import (
	"log"
	"os"
	"os/user"

	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/crypto/bcrypt"
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
		os.Mkdir(cp, os.ModePerm)
	}

	return cp
}

func InitConfig() error {
	cp := GenConfigPath()
	log.Println(cp)
	var err error
	ConfigDb, err = leveldb.OpenFile(cp, nil)
	if err != nil {
		return err
	}

	return nil
}

func NewUsr(root string, p string) error {
	usr := GetCurrentUsr()
	if root == "" {
		root = usr.HomeDir
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	m := make(map[string]string)
	m["password"] = string(hash)
	m["root"] = root

	return BatchPut(m)
}

func BatchPut(m map[string]string) error {
	b := new(leveldb.Batch)
	for k, v := range m {
		b.Put([]byte(k), []byte(v))
	}

	return ConfigDb.Write(b, nil)
}
