package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/orbit-drive/orbit-drive/utils"
)

const (
	CONFIGFILENAME string = "config.json"
)

var (
	cfg *Config
	// ErrSecretPhraseNotProvided is returned when initializing a config with no secrete phrase
	ErrSecretPhraseNotProvided = errors.New("config: no secret phrase provided")
)

// Config represents the usr configuration settings
type Config struct {
	// Root is the absolute path of the directory to synchronize.
	Root string `json:"root_path"`

	// SecretPhrase is the user authentication token for synchronization.
	SecretPhrase string `json:"secret_phrase"`

	// NodeAddr is address of the ipfs node for the api request. (Default: infura)
	NodeAddr string `json:"node_addr"`

	// Port to use by p2p connections.
	P2PPort string `json:"p2p_port"`
}

func GetNodeAddr() string {
	return cfg.NodeAddr
}

func GetRoot() string {
	return cfg.Root
}

// NewConfig initialize a new usr config and save it to config file.
func NewConfig(root, secretPhrase, nodeAddr, p2pPort string) error {
	if secretPhrase == "" {
		return ErrSecretPhraseNotProvided
	}
	configFile, err := createConfigFile()
	if err != nil {
		return err
	}
	defer configFile.Close()

	spHash, err := utils.SecureHash(secretPhrase)
	if err != nil {
		return err
	}

	config := &Config{
		Root:         root,
		SecretPhrase: string(spHash),
		NodeAddr:     nodeAddr,
		P2PPort:      p2pPort,
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	_, err = configFile.Write(configData)
	return err
}

// LoadConfig reads config from config.json file.
func LoadConfig(nodeAddr, p2pPort string) error {
	configPath := configFilePath()
	cfg = &Config{}
	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	parser := json.NewDecoder(configFile)
	if err = parser.Decode(cfg); err != nil {
		return err
	}
	if nodeAddr != "" {
		cfg.NodeAddr = nodeAddr
	}
	if p2pPort != "" {
		cfg.P2PPort = p2pPort
	}
	return nil
}

func createConfigFile() (*os.File, error) {
	configFilePath := configFilePath()
	if utils.PathExists(configFilePath) {
		f, err := os.Open(configFilePath)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	return os.Create(configFilePath)
}

func configFilePath() string {
	configDir := utils.GetConfigDir()
	return filepath.Join(configDir, CONFIGFILENAME)
}
