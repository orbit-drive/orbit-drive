package fs

import (
	"crypto/sha256"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// Dir helpers
func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

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

// String parsing helpers
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
