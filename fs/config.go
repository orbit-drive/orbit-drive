package fs

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

const (
	CONFIG_KEY string = "CONFIG"
)

// Config represents the usr configuration settings
type Config struct {
	// Root is the absolute path of the directory to synchronize.
	Root string `json:"root_path"`

	// Node is address of the ipfs node for the api request. (Default: infura)
	Node string `json:"node_addr"`

	// Password is the usr password set used for file encryption.
	Password string `json:"password_hash"`
}

// NewConfig initialize a new usr config and save it.
func NewConfig(root, node, p string) error {
	if root == "" {
		root = GetCurrentDir()
	}

	hash, err := bcrypt.GenerateFromPassword(ToByte(p), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	c := &Config{
		Root:     root,
		Node:     node,
		Password: ToStr(hash),
	}
	return c.Save()
}

// LoadConfig loads a stored config from: (defaults: ~/.orbit-drive/.config)
func LoadConfig() (*Config, error) {
	data, err := Db.Get(ToByte(CONFIG_KEY), nil)
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
	return Db.Put(ToByte(CONFIG_KEY), p, nil)
}
