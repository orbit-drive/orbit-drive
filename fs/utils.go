package fs

import (
	"crypto/sha256"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func HashStr(p string) []byte {
	hash := sha256.Sum256(ToByte(p))
	return hash[:]
}

func ToByte(s string) []byte {
	return []byte(s)
}

func ToStr(b []byte) string {
	return string(b[:])
}
