package sync

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/kdisneur/dropbox_sync/pkg/dropbox"
	"github.com/kdisneur/dropbox_sync/pkg/local"
	"github.com/sirupsen/logrus"
)

type Sync struct {
	Client         *dropbox.Client
	DropboxLogger  *logrus.Entry
	DropboxScanner *dropbox.Scanner
	LocalScanner   *local.Scanner
	LocalBasePath  string
	LocalLogger    *logrus.Entry
	RemoteBasePath string
}

// NewSync creates a new bidirectional synchronizer between dropbox and the local filesystem
func NewSync(client *dropbox.Client, localPath string, remotePath string) *Sync {
	dropboxLogger := logrus.WithFields(
		logrus.Fields{"folder": remotePath, "direction": "dropbox-to-local"},
	)

	localLogger := logrus.WithFields(
		logrus.Fields{"folder": localPath, "direction": "local-to-dropbox"},
	)

	return &Sync{
		Client:         client,
		LocalBasePath:  localPath,
		RemoteBasePath: remotePath,
		DropboxScanner: dropbox.NewScanner(dropboxLogger, *client, remotePath),
		DropboxLogger:  dropboxLogger,
		LocalScanner:   local.NewScanner(localLogger, localPath),
		LocalLogger:    localLogger,
	}
}

// DropboxFolder copies dropbox files to a local folder
func (s *Sync) DropboxFolder() error {
	for s.DropboxScanner.Next() {
		action := s.DropboxScanner.Entry()

		var err error
		switch action.Type {
		case dropbox.ActionTypeCreate:
			s.DropboxLogger.Debugf("creates or update file or folder '%s'", action.File.RelativePath)
			err = s.createLocalFileOrFolder(action.File)
			s.LocalScanner.NotifyCreation(action.File.RelativePath)
		case dropbox.ActionTypeDelete:
			s.DropboxLogger.Debugf("delete file or folder '%s'", action.File.RelativePath)
			err = s.deleteLocalFileOrFolder(action.File)
			s.LocalScanner.NotifyDeletion(action.File.RelativePath)
		default:
			err = fmt.Errorf("unsupported dropbox action: %s", action.Type)
		}

		if err != nil {
			return err
		}
	}

	if s.DropboxScanner.Err() != nil {
		return s.DropboxScanner.Err()
	}

	return nil
}

// LocalFolder copies dropbox files to a local folder
func (s *Sync) LocalFolder() error {
	for s.LocalScanner.Next() {
		action := s.LocalScanner.Entry()

		var err error
		switch action.Type {
		case local.ActionTypeCreate:
			s.LocalLogger.Debugf("creates or update file or folder '%s'", action.File.RelativePath)
			err = s.createDropboxFileOrFolder(action.File)
		case local.ActionTypeDelete:
			s.LocalLogger.Debugf("delete file or folder '%s'", action.File.RelativePath)
			err = dropbox.FileDelete(*s.Client, path.Join(s.RemoteBasePath, action.File.RelativePath))
		default:
			err = fmt.Errorf("unsupported local action: %s", action.Type)
		}

		if err != nil {
			return err
		}
	}

	if s.LocalScanner.Err() != nil {
		return s.LocalScanner.Err()
	}

	return nil
}

func (s *Sync) createDropboxFileOrFolder(file local.File) error {
	switch file.Type {
	case local.FileTypeFile:
		content, err := ioutil.ReadFile(file.Path)
		if err != nil {
			return err
		}

		return dropbox.FileUpload(*s.Client, path.Join(s.RemoteBasePath, file.RelativePath), content)
	case local.FileTypeFolder:
		return dropbox.FolderCreate(*s.Client, path.Join(s.RemoteBasePath, file.RelativePath))
	default:
		return fmt.Errorf("unsupported local file type: %s", file.Type)
	}
}

func (s *Sync) createLocalFileOrFolder(file dropbox.File) error {
	filePath := path.Join(s.LocalBasePath, file.RelativePath)

	currentSum, currentSumErr := dropbox.HashFromFile(filePath)
	if currentSumErr == nil && currentSum == file.ContentHash {
		s.DropboxLogger.Debugf("file already up-to-date. skip creation (%s)", filePath)
		return nil
	}

	switch file.Type {
	case dropbox.FileTypeFolder:
		return os.MkdirAll(filePath, 0750)
	case dropbox.FileTypeFile:
		return s.fetchDropboxContent(file.RemotePath, filePath)
	default:
		return fmt.Errorf("unsupported dropbox file type: %s", file.Type)
	}
}

func (s *Sync) deleteLocalFileOrFolder(file dropbox.File) error {
	filePath := path.Join(s.LocalBasePath, file.RelativePath)

	return os.RemoveAll(filePath)
}

func (s *Sync) fetchDropboxContent(dropboxPath string, localPath string) error {
	content, err := dropbox.FileDownload(*s.Client, dropboxPath)
	if err != nil {
		return err
	}

	err = writeLocalFileAndSubfolders(localPath, content)
	if err != nil {
		return err
	}

	return nil
}

func writeLocalFileAndSubfolders(filePath string, content []byte) error {
	err := os.MkdirAll(path.Dir(filePath), 0750)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, content, 0640)
	if err != nil {
		return err
	}

	return nil
}
