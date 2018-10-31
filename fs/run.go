package fs

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wlwanpan/orbit-drive/api"
	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/db"
)

func Run(c *db.Config) {
	fmt.Println("Starting orbit file sync, watching: ", c.Root)
	defer fmt.Println("Orbit sync stopped.")

	// Init ipfs api shell
	api.InitShell(c.Node)

	// Get previouly stored files.
	fStore, err := db.GetFileStore()
	if err != nil {
		log.Println(err)
	}

	// Load Tree from Db and Gen diffing Tree
	err = InitVTree(c.Root, fStore)
	if err != nil {
		// Delete prev files saved but no longer present in file system.
		fStore.Dump()
	}

	// Logs the json representation of the loaded VTree
	data, err := json.MarshalIndent(&VTree, "", "	")
	if err != nil {
		log.Println(err)
	}
	log.Println(common.ToStr(data))

	// Init and Start file watcher
	w := NewWatcher(c.Root)
	w.Start()
}
