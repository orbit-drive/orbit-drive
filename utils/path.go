package utils

import (
	"os"
	"os/user"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	CONFIGPATH string = "/.orbit-drive"
)

func ExtractFileName(path string) string {
	return filepath.Base(path)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func IsHidden(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Warn(err)
		return false
	}
	filename := fi.Name()
	return filename[0:1] == "."
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

func GetConfigDir() string {
	return GetHomeDir() + CONFIGPATH
}
