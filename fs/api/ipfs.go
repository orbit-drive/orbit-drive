package api

import (
	"log"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

var (
	// Shell holds a ipfs shell instance for access to the
	// the ipfs network. (Default: Infura node)
	Shell *shell.Shell
)

// InitShell initialize a new shell.
func InitShell(addr string) {
	Shell = shell.NewShell(addr)
}

// UploadFile takes a file path and upload it to ipfs
// and return the generate hash.
func UploadFile(p string) (string, error) {
	f, _ := os.Open(p)
	log.Println("Uploading file: ", p)
	cid, err := Shell.Add(f)
	if err != nil {
		return "", err
	}

	log.Println("Uploaded file: ", cid)
	return cid, nil
}
