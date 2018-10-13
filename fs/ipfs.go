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

func Upload(p string) error {
	fi, err := os.Stat(p)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		log.Println("Uploading dir: ", p)
		cid, err := Shell.AddDir(p)
		if err != nil {
			return err
		}
		log.Println("Uploaded dir: ", cid)
		return nil
	}

	f, _ := os.Open(p)
	log.Println("Uploading file: ", p)
	cid, err := Shell.Add(f)
	if err != nil {
		return err
	}

	log.Println("Uploaded file: ", cid)

	return nil
}
