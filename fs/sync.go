package fs

import (
	"encoding/json"
	"fmt"
	"log"
)

func Sync(c *Config) {
	fmt.Println("Starting orbit file sync, watching: ", c.Root)

	// Init ipfs api shell
	InitShell(c.Node)

	// Load Tree from Db and Gen diffing Tree
	InitTree()
	t, err := GenTreeFromPath(c)
	if err != nil {
		log.Fatal(err)
	}

	// Diff t and Tree
	data, _ := json.Marshal(t)
	log.Println(string(data[:]))

	// Init and Start file watcher
	// w := NewWatcher(c.Root)
	// w.Start()
}
