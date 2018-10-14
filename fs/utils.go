package fs

import (
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
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(path)
}
