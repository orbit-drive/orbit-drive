package fs

import (
	"encoding/json"
	"fmt"
	"log"
)

func Sync(c *Config) {
	fmt.Println("Starting orbit file sync, watching: ", c.Root)
	defer fmt.Println("Orbit sync stopped.")

	// Init ipfs api shell
	InitShell(c.Node)

	// Get previouly stored files.
	fStore, err := GetFileStore()
	if err != nil {
		log.Println(err)
	}

	// Load Tree from Db and Gen diffing Tree
	err = InitVTree(c.Root, fStore)
	if err != nil {
		// Delete prev files saved but no longer present in file system.
		RunGarbageCollection(fStore)
	}

	// Logs the json representation of the loaded VTree
	data, err := json.Marshal(&VTree)
	if err != nil {
		log.Println(err)
	}
	log.Println(ToStr(data))

	// Init and Start file watcher
	w := NewWatcher(c.Root)
	w.Start()
}
