package fs

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

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

// NewConfig initialize a new usr config and save it to config file.
func NewConfig(root, secretPhrase, nodeAddr string) error {
	if secretPhrase == "" {
		return errors.New("no secret phrase provided")
	}
	configFile, err := createConfigFile()
	if err != nil {
		return nil
	}
	defer configFile.Close()

	spHash, err := fsutil.SecureHash(secretPhrase)
	if err != nil {
		return err
	}

	config := &Config{
		Root:         root,
		SecretPhrase: string(spHash),
		NodeAddr:     nodeAddr,
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	_, err = configFile.Write(configData)
	return err
}

// LoadConfig reads config from config.json file.
func LoadConfig() (*Config, error) {
	configPath := configFilePath()
	config := &Config{}
	configFile, err := os.Open(configPath)
	if err != nil {
		return &Config{}, err
	}
	parser := json.NewDecoder(configFile)
	if err = parser.Decode(config); err != nil {
		return &Config{}, err
	}
	return config, nil
}

func createConfigFile() (*os.File, error) {
	configFilePath := configFilePath()
	if fsutil.PathExists(configFilePath) {
		return &os.File{}, nil
	}
	return os.Create(configFilePath)
}

func configFilePath() string {
	configDir := fsutil.GetConfigDir()
	return filepath.Join(configDir, "config.json")
}
