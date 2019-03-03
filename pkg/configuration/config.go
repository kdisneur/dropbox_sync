package configuration

import (
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

// Config represents the configuration file
type Config struct {
	Authentication DropboxAuthentication `toml:"authentication"`
	Folders        []Folder              `toml:"folder"`
}

// DropboxAuthentication represents the Dropbox authentication configuration
type DropboxAuthentication struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}

// Folder represents a folder to synchronize
type Folder struct {
	RemotePath string `toml:"remote_path"`
	LocalPath  string `toml:"local_path"`
}

// LoadConfiguration load the configuration from the home folder
func LoadConfiguration() (*Config, error) {
	filePath, err := homedir.Expand(path.Join("~", ".config", "dropbox_sync", "config"))
	if err != nil {
		return nil, errors.Wrap(err, "can't find HOME folder")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "can't read config file")
	}

	rawConfig, err := toml.LoadReader(file)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse TOML config file")
	}

	config := &Config{}
	err = rawConfig.Unmarshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse TOML config file")
	}

	for i, folder := range config.Folders {
		localPath, err := homedir.Expand(folder.LocalPath)
		if err != nil {
			return nil, errors.Wrap(err, "can't expand local path")
		}
		config.Folders[i].LocalPath = localPath
	}

	return config, nil
}
