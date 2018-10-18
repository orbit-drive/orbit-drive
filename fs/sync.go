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

	// Init saved files
	InitSavedFiles()

	// Load Tree from Db and Gen diffing Tree
	err := InitVTree(c.Root)
	if err != nil {
		// Delete prev files saved but no longer present in file system.
		RunGarbageCollect()
	}

	// Logs the json representation of the loaded VTree
	data, err := json.Marshal(&VTree)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(data[:]))

	// Diff t and Tree
	// Init and Start file watcher
	// w := NewWatcher(c.Root)
	// w.Start()
}
