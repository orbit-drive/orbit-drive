package main

import (
	"flag"
	"log"
	"os"

	"github.com/wlwanpan/ipfs-drive/fs"
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

	ipfsync := fs.NewIpfsSync(*rootPath, *nodeAddr)

	iter := ipfsync.Db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		log.Printf("k: %s, v: %s", key, value)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Fatal(err)
	}

	ipfsync.Start()
}
