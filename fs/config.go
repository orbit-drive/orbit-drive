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

	// AuthToken is the user authentication token for synchronization.
	AuthToken string `json:"auth_token"`

	// NodeAddr is address of the ipfs node for the api request. (Default: infura)
	NodeAddr string `json:"node_addr"`

	// HubAddr is the address of the backend service for device sync.
	HubAddr string `json:"hub_addr"`

	// Password is the usr password set used for file encryption.
	Password string `json:"password_hash"`
}

// NewConfig initialize a new usr config and save it.
func NewConfig(root, authToken, nodeAddr, hubAddr, path string) error {
	if root == "" {
		root = fsutil.GetCurrentDir()
	}

	hash, err := fsutil.PasswordHash(path)
	if err != nil {
		return err
	}

	c := &Config{
		Root:      root,
		AuthToken: authToken,
		NodeAddr:  nodeAddr,
		HubAddr:   hubAddr,
		Password:  fsutil.ToStr(hash),
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

// Update overwrites the config hub and ipfs address if a non zero value is provided.
func (c *Config) Update(n string, h string) {
	if n != "" {
		c.NodeAddr = n
	}
	if h != "" {
		c.HubAddr = h
	}
}

// Save persist the current configuration.
func (c *Config) Save() error {
	p, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return db.Put(fsutil.ToByte(fsutil.CONFIGKEY), p)
}
