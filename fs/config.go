package fs

import (
	"encoding/json"

	"github.com/orbit-drive/orbit-drive/fs/db"
	"github.com/orbit-drive/orbit-drive/fsutil"
)

// Config represents the usr configuration settings
type Config struct {
	// Root is the absolute path of the directory to synchronize.
	Root string `json:"root_path"`

	// SecretPhrase is the user authentication token for synchronization.
	SecretPhrase string `json:"secret_phrase"`

	// NodeAddr is address of the ipfs node for the api request. (Default: infura)
	NodeAddr string `json:"node_addr"`
}

// ---------------------------------------------------------
// Refact config to using config file / Dont save in leveldb
// https://micro.mu/docs/go-config.html

// NewConfig initialize a new usr config and save it.
func NewConfig(root, secretPhrase, nodeAddr string) error {
	if root == "" {
		root = fsutil.GetCurrentDir()
	}

	spHash, err := fsutil.SecureHash(secretPhrase)
	if err != nil {
		return err
	}

	c := &Config{
		Root:         root,
		SecretPhrase: string(spHash),
		NodeAddr:     nodeAddr,
	}
	return c.Save()
}

// LoadConfig loads a stored config from: (defaults: ~/.orbit-drive/.config)
func LoadConfig() (*Config, error) {
	data, err := db.Get(fsutil.ToByte(fsutil.CONFIGKEY))
	if err != nil {
		return &Config{}, err
	}
	c := &Config{}
	err = json.Unmarshal(data, c)
	if err != nil {
		return &Config{}, err
	}
	return c, nil
}

// Save persist the current configuration.
func (c *Config) Save() error {
	p, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return db.Put(fsutil.ToByte(fsutil.CONFIGKEY), p)
}
