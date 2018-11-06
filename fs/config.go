package fs

import (
	"encoding/json"

	"github.com/wlwanpan/orbit-drive/common"
	"github.com/wlwanpan/orbit-drive/fs/db"
)

// Config represents the usr configuration settings
type Config struct {
	// Root is the absolute path of the directory to synchronize.
	Root string `json:"root_path"`

	// NodeAddr is address of the ipfs node for the api request. (Default: infura)
	NodeAddr string `json:"node_addr"`

	// HubAddr is the address of the backend service for device sync.
	HubAddr string `json:"hub_addr"`

	// Password is the usr password set used for file encryption.
	Password string `json:"password_hash"`
}

// NewConfig initialize a new usr config and save it.
func NewConfig(root, nodeAddr, p string) error {
	if root == "" {
		root = common.GetCurrentDir()
	}

	hash, err := common.PasswordHash(p)
	if err != nil {
		return err
	}

	c := &Config{
		Root:     root,
		NodeAddr: nodeAddr,
		Password: common.ToStr(hash),
	}
	return c.Save()
}

// LoadConfig loads a stored config from: (defaults: ~/.orbit-drive/.config)
func LoadConfig() (*Config, error) {
	data, err := db.Get(common.ToByte(common.CONFIG_KEY))
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
	return db.Put(common.ToByte(common.CONFIG_KEY), p)
}
