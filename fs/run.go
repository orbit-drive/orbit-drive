package fs

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/fs/api"
	"github.com/wlwanpan/orbit-drive/fs/db"
	"github.com/wlwanpan/orbit-drive/fs/vtree"
)

func Run(c *Config) {
	fmt.Println("Starting orbit file sync, watching: ", c.Root)
	defer fmt.Println("Orbit sync stopped.")

	// Init ipfs api shell
	api.InitShell(c.Node)

	// Get previouly stored files.
	sources, err := db.GetSources()
	if err != nil {
		log.Println(err)
	}

	// Load Tree from Db and Gen diffing Tree
	err = vtree.InitVTree(c.Root, sources)
	if err != nil {
		// Delete prev files saved but no longer present in file system.
		sources.Dump()
	}

	// Logs the json representation of the loaded VTree
	data, err := json.MarshalIndent(&vtree.VTree, "", "	")
	if err != nil {
		log.Println(err)
	}
	log.Println(common.ToStr(data))

	// Init and Start file watcher
	w := NewWatcher(c.Root)
	w.Start()
}
