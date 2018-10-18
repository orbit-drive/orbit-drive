package fs

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

const (
	CONFIG_KEY string = "CONFIG"
)

type Config struct {
	Root     string `json:"root_path"`
	Node     string `json:"node_addr"`
	Password string `json:"password_hash"`
}

func (c *Config) Save() error {
	p, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return Db.Put([]byte(CONFIG_KEY), p, nil)
}

func LoadConfig() (*Config, error) {
	data, err := Db.Get([]byte(CONFIG_KEY), nil)
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

func NewConfig(root, node, p string) error {
	if root == "" {
		root = GetHomeDir()
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	c := &Config{
		Root:     root,
		Node:     node,
		Password: string(hash[:]),
	}
	return c.Save()
}
