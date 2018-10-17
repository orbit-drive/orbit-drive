package fs

import (
	"log"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

var (
	Shell *shell.Shell
)

func InitShell(addr string) {
	Shell = shell.NewShell(addr)
}

func UploadFile(p string) (string, error) {
	fi, err := os.Stat(p)
	if err != nil {
		return "", err
	}
	if fi.IsDir() {
		log.Println("Path provided is not a file: ", p)
		return "", ErrInValidPath
	}

	f, _ := os.Open(p)
	log.Println("Uploading file: ", p)
	cid, err := Shell.Add(f)
	if err != nil {
		return "", err
	}

	log.Println("Uploaded file: ", cid)
	return cid, nil
}
