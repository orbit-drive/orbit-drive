package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/wlwanpan/orbit-drive/fs"
)

func main() {
	p := argparse.NewParser("orbit-drive", "File uploader and synchronizer built on IPFS and Infura.")

	// init command
	initCmd := p.NewCommand("init", "Initialize folder to sync.")
	root := initCmd.String("r", "root", &argparse.Options{
		Required: false,
		Help:     "Root path of folder to synchronise.",
	})
	nodeAddr := initCmd.String("n", "node", &argparse.Options{
		Required: false,
		Default:  "https://ipfs.infura.io:5001",
		Help:     "Ipfs node address, will default to an infura node if none is provided.",
	})
	password := initCmd.String("p", "password", &argparse.Options{
		Required: false,
		Help:     "Set password for file encryption.",
	})

	// sync command
	syncCmd := p.NewCommand("sync", "Start syncing folder to the ipfs network.")

	err := p.Parse(os.Args)
	if err != nil {
		fmt.Println(p.Usage(err))
		os.Exit(0)
	}

	fs.InitDb()
	defer fs.Db.Close()

	switch true {
	case initCmd.Happened():
		err := fs.NewConfig(*root, *nodeAddr, *password)
		if err != nil {
			fmt.Println(p.Usage(err))
			os.Exit(0)
		}
		fmt.Println("Configured! Run the following command to start syncing: orbit-drive sync")
	case syncCmd.Happened():
		c, err := fs.LoadConfig()
		if err != nil {
			fmt.Println(p.Usage(err))
			os.Exit(0)
		}
		fs.Sync(c)
	default:
		os.Exit(0)
	}
}
