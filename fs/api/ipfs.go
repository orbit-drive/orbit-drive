package api

import (
	"os"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/wlwanpan/orbit-drive/fs/sys"
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
	file, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sys.Notify("Uploading: ", file.Name())
	cid, err := Shell.Add(file)
	if err != nil {
		return "", err
	}

	sys.Notify("Uploaded: ", cid)
	return cid, nil
}
