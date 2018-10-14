package fs

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

const (
	configKey string = "config"
)

type Config struct {
	Root     string `json:"root_path"`
	Password string `json:"password_hash"`
}

func (c *Config) Save() error {
	p, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return Db.Put([]byte(configKey), p, nil)
}

func (c *Config) Load() error {
	p, err := Db.Get([]byte(configKey), nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(p, c)
}

func NewConfig(root string, p string) error {
	if root == "" {
		root = GetHomeDir()
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	c := &Config{
		Root:     root,
		Password: string(hash[:]),
	}
	return c.Save()
}
