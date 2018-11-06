package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/wlwanpan/orbit-drive/fs"
	"github.com/wlwanpan/orbit-drive/fs/db"
)

func main() {
	p := argparse.NewParser("orbit-drive", "File uploader and synchronizer built on IPFS and Infura.")

	// init command
	initCmd := p.NewCommand("init", "Initialize folder to sync.")
	root := initCmd.String("r", "root", &argparse.Options{
		Required: false,
		Help:     "Root path of folder to synchronise.",
	})
	nodeAddr := initCmd.String("n", "node-addr", &argparse.Options{
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
	hubAddr := syncCmd.String("b", "hub-addr", &argparse.Options{
		Required: false,
		Default:  "",
		Help:     "Hub address for device synchronization.",
	})

	if err := p.Parse(os.Args); err != nil {
		log.Fatal(p.Usage(err))
	}

	db.InitDb()
	defer db.CloseDb()

	switch {
	case initCmd.Happened():
		err := fs.NewConfig(*root, *nodeAddr, *password)
		if err != nil {
			log.Fatal(p.Usage(err))
		}
		fmt.Println("Configured! Run the following command to start syncing: orbit-drive sync")
	case syncCmd.Happened():
		c, err := fs.LoadConfig()
		if err != nil {
			log.Fatal(p.Usage(err))
		}
		c.Update("", *hubAddr)
		fs.Run(c)
	default:
		os.Exit(0)
	}
}
