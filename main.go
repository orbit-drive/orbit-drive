package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/wlwanpan/ip-drive/fs"
)

const (
	InfuraAddr = "https://ipfs.infura.io:5001"
)

func main() {
	p := argparse.NewParser("ip-drive", "File uploader and synchronizer built on IPFS and Infura.")

	// init command
	initCmd := p.NewCommand("init", "Initialize folder to sync.")
	root := initCmd.String("r", "root", &argparse.Options{
		Required: false,
		Help:     "Root path of folder to synchronise.",
	})
	password := initCmd.String("p", "password", &argparse.Options{
		Required: false,
		Help:     "Set password for file encryption.",
	})

	// sync command
	syncCmd := p.NewCommand("sync", "Start syncing folder to the ipfs network.")
	nodeAddr := syncCmd.String("n", "node", &argparse.Options{
		Required: false,
		Default:  InfuraAddr,
		Help:     "Ipfs node address, will default to an infura node if none is provided.",
	})

	err := p.Parse(os.Args)
	if err != nil {
		fmt.Println(p.Usage(err))
		os.Exit(1)
	}

	fs.InitConfig()
	defer fs.ConfigDb.Close()

	switch true {
	case initCmd.Happened():
		fs.NewUsr(*root, *password)
	case syncCmd.Happened():
		fs.InitShell(*nodeAddr)
		ipfsync := fs.NewSync(*root, *nodeAddr)
		ipfsync.Start()
	default:
		os.Exit(1)
	}
}
