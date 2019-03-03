package configuration

import (
	"io/ioutil"
	"path"

	"github.com/kdisneur/dropbox_sync/pkg/dropbox"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

var tokenFilePath = path.Join("~", ".config", "dropbox_sync", "token")

// LoadDropboxClient load the dropbox client from a stored token
func LoadDropboxClient() (*dropbox.Client, error) {
	filePath, err := homedir.Expand(tokenFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "can't find HOME folder")
	}

	token, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "can't read token file")
	}

	client := dropbox.NewClient(string(token))

	return &client, nil
}

// SaveDropboxToken store the token on file and returns a client
func SaveDropboxToken(token string) (*dropbox.Client, error) {
	filePath, err := homedir.Expand(tokenFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "can't find HOME folder")
	}

	err = ioutil.WriteFile(filePath, []byte(token), 0600)
	if err != nil {
		return nil, errors.Wrap(err, "can't write token file")
	}

	client := dropbox.NewClient(token)

	return &client, nil
}
