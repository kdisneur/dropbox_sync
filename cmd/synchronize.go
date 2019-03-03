package cmd

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"

	"github.com/kdisneur/dropbox_sync/pkg/configuration"
	"github.com/kdisneur/dropbox_sync/pkg/dropbox"
	"github.com/kdisneur/dropbox_sync/pkg/sync"
	"github.com/sirupsen/logrus"
)

// Synchronize synchronize data between Dropbox and a local folder
type Synchronize struct{}

// Run starts the Dropbox <-> folder synchronization
func (s Synchronize) Run() {
	var client *dropbox.Client
	var err error

	config, err := configuration.LoadConfiguration()
	if err != nil {
		fail(err)
	}

	client, err = configuration.LoadDropboxClient()
	if err != nil {
		client, err = s.authenticate(config)
	}

	if err != nil {
		fail(err)
	}

	waitingErrors := make(chan error, 0)

	for _, folder := range config.Folders {
		os.MkdirAll(folder.LocalPath, 0755)

		synchronizer := sync.NewSync(client, folder.LocalPath, folder.RemotePath)

		go s.startScanningDropbox(synchronizer, waitingErrors)
		go s.startScanningLocal(synchronizer, waitingErrors)
	}

	err = <-waitingErrors
	fail(err)
}

func (s Synchronize) authenticate(config *configuration.Config) (*dropbox.Client, error) {
	oauth2 := dropbox.NewOAuth2(config.Authentication.ClientID, config.Authentication.ClientSecret)

	fmt.Println("dropbox token not found. starts the authentication process.")
	fmt.Printf("open your browser to authenticate: %s\n", oauth2.AuthorizationURL())
	fmt.Printf("enter the code: ")
	authorizationCode, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println()
	if err != nil {
		return nil, err
	}

	token, err := oauth2.GetAccessToken(string(authorizationCode))
	if err != nil {
		return nil, err
	}

	client, err := configuration.SaveDropboxToken(token)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s Synchronize) startScanningDropbox(synchronizer *sync.Sync, errors chan error) {
	logrus.Infof("start syncing Dropbox folder '%s' to local '%s' path", synchronizer.RemoteBasePath, synchronizer.LocalBasePath)
	err := synchronizer.DropboxFolder()
	if err != nil {
		errors <- err
	}
}

func (s Synchronize) startScanningLocal(synchronizer *sync.Sync, errors chan error) {
	logrus.Infof("start syncing local folder '%s' to Dropbox '%s' path", synchronizer.LocalBasePath, synchronizer.RemoteBasePath)
	err := synchronizer.LocalFolder()
	if err != nil {
		errors <- err
	}
}
