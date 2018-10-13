package main

import (
	"flag"
	"log"
	"os"

	"github.com/wlwanpan/ip-drive/fs"
)

const (
	infuraAddr = "https://ipfs.infura.io:5001"
)

func main() {

	nodeAddr := flag.String("ipfs-addr", infuraAddr, "Ipfs node address.")
	rootPath := flag.String("root", "", "Root path of folder to synchronize.")

	flag.Parse()

	if *rootPath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		*rootPath = pwd
	}

	log.Println("Syncing: ", *rootPath)

	// Create a dir tree mapper to diff sync
	// files, err := ioutil.ReadDir(*rootPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, f := range files {
	// 	log.Println(f.Name())
	// }

	fs.InitShell(*nodeAddr)
	fs.InitConfig()
	defer fs.ConfigDb.Close()

	ipfsync := fs.NewSync(*rootPath, *nodeAddr)
	ipfsync.Start()
}
