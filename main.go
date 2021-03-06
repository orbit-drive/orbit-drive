package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
	"github.com/orbit-drive/orbit-drive/config"
	"github.com/orbit-drive/orbit-drive/db"
	"github.com/orbit-drive/orbit-drive/sync"
	"github.com/orbit-drive/orbit-drive/utils"
	log "github.com/sirupsen/logrus"
)

func initLogger() *os.File {
	logFilePath := filepath.Join(utils.GetConfigDir(), "info.log")
	if !utils.PathExists(logFilePath) {
		os.Create(logFilePath)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	return logFile
}

func main() {
	p := argparse.NewParser("orbit-drive", "File uploader and synchronizer built on IPFS and Infura.")

	// init command
	initCmd := p.NewCommand("init", "Initialize folder to sync.")
	root := initCmd.String("r", "root", &argparse.Options{
		Required: false,
		Default:  utils.GetCurrentDir(),
		Help:     "Root path of folder to synchronise.",
	})
	secretPhrase := initCmd.String("s", "secret", &argparse.Options{
		Required: true,
		Default:  "",
		Help:     "Set a secret phrase and share with our devices you with to sync with.",
	})

	// sync command
	syncCmd := p.NewCommand("sync", "Start syncing folder to the ipfs network.")

	// Optional command
	nodeAddr := p.String("n", "node-addr", &argparse.Options{
		Required: false,
		// Default:  "https://ipfs.infura.io:5001",
		Help: "Ipfs node address, will default to an infura node if none is provided.",
	})
	p2pPort := p.String("p", "p2p-port", &argparse.Options{
		Required: false,
		// Default:  "6666",
		Help: "P2P port to use to connect to other peers via tcp.",
	})

	// TODO: Add check if port is in use.
	if err := p.Parse(os.Args); err != nil {
		log.Fatal(p.Usage(err))
	}

	if err := db.InitDb(); err != nil {
		log.Fatal(err)
	}
	defer db.CloseDb()

	switch {
	case initCmd.Happened():
		err := config.NewConfig(*root, *secretPhrase, *nodeAddr, *p2pPort)
		if err != nil {
			log.Fatal(p.Usage(err))
		}
		fmt.Println("Configured! Run the following command to start syncing: orbit-drive sync")
	case syncCmd.Happened():
		c, err := config.LoadConfig(*nodeAddr, *p2pPort)
		if err != nil {
			log.Fatal(err)
		}

		f := initLogger()
		defer f.Close()
		sync.Run(c)
	default:
		os.Exit(0)
	}
}
