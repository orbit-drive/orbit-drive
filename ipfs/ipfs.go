package ipfs

import (
	"errors"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/orbit-drive/orbit-drive/sys"
)

var (
	// ErrNodeOffline is returned when performing an operation to a
	// disconnected ipfs node.
	ErrNodeOffline = errors.New("ipfs: node not live")

	// ErrNodeNotInitialized is returned when accessing a nil pointer to
	// the ipfs shell instance.
	ErrNodeNotInitialized = errors.New("ipfs: node not initialized")
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

// IsLive return true if ipfs node is live.
func IsLive() (bool, error) {
	if Shell == nil {
		return false, ErrNodeNotInitialized
	}
	return Shell.IsUp(), nil
}

// UploadFile takes a file path and upload it to ipfs
// and return the generate hash.
func UploadFile(p string) (string, error) {
	isLive, err := IsLive()
	if err != nil {
		return "", err
	}
	if !isLive {
		return "", ErrNodeOffline
	}
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
