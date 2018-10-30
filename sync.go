package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/db"
	"github.com/wlwanpan/orbit-drive/fs"
)

func Sync(c *db.Config) {
	fmt.Println("Starting orbit file sync, watching: ", c.Root)
	defer fmt.Println("Orbit sync stopped.")

	// Init ipfs api shell
	fs.InitShell(c.Node)

	// Get previouly stored files.
	fStore, err := db.GetFileStore()
	if err != nil {
		log.Println(err)
	}

	// Load Tree from Db and Gen diffing Tree
	err = fs.InitVTree(c.Root, fStore)
	if err != nil {
		// Delete prev files saved but no longer present in file system.
		fStore.Dump()
	}

	// Logs the json representation of the loaded VTree
	data, err := json.MarshalIndent(&fs.VTree, "", "	")
	if err != nil {
		log.Println(err)
	}
	log.Println(common.ToStr(data))

	// Init and Start file watcher
	w := fs.NewWatcher(c.Root)
	w.Start()
}
