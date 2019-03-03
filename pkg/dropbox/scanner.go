package dropbox

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/kdisneur/dropbox_sync/pkg/dropbox/internal"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Scanner represents a paginated list of entries
type Scanner struct {
	Client      Client
	buffer      []Action
	err         error
	hasNextPage *bool
	index       int
	mutex       sync.Mutex
	nextCursor  string
	path        string
	logger      *logrus.Entry
}

// NewScanner creates a new folder scanner
func NewScanner(logger *logrus.Entry, client Client, path string) *Scanner {
	return &Scanner{Client: client, logger: logger, path: path}
}

// Next replace the `Entry` with the following one if it can and return false if it can't
func (f *Scanner) Next() bool {
	if f.Err() != nil {
		return false
	}

	if f.buffer == nil {
		return f.loadFirstPage()
	}

	nextIndex := f.index + 1
	if nextIndex < len(f.buffer) {
		f.index = nextIndex
		return true
	}

	if f.hasNextPage != nil && !*f.hasNextPage {
		err := f.waitForUpdate()
		if err != nil {
			f.err = err
			return false
		}
	}

	if !f.loadNextPage() {
		return f.Next()
	}

	return true
}

// Entry returns the current scanner content
func (f *Scanner) Entry() *Action {
	if f.Err() != nil {
		return nil
	}

	if f.index >= len(f.buffer) {
		return nil
	}

	return &f.buffer[f.index]
}

// Err returns the scanner error if one exists
func (f *Scanner) Err() error {
	return f.err
}

func (f *Scanner) loadFirstPage() bool {
	f.logger.WithFields(logrus.Fields{"cursor": f.nextCursor}).Debugf("fetch first page of entries")

	return f.executeQuery(func() ([]byte, error) {
		return internal.POSTWithBody(
			"https://api.dropboxapi.com/2/files/list_folder",
			f.Client.token,
			map[string]interface{}{
				"path":                    f.path,
				"recursive":               true,
				"include_media_info":      false,
				"include_deleted":         false,
				"include_mounted_folders": true,
			},
		)
	})
}

func (f *Scanner) loadNextPage() bool {
	f.logger.WithFields(logrus.Fields{"cursor": f.nextCursor}).Debugf("fetch next page of entries")

	return f.executeQuery(func() ([]byte, error) {
		return internal.POSTWithBody(
			"https://api.dropboxapi.com/2/files/list_folder/continue",
			f.Client.token,
			map[string]interface{}{"cursor": f.nextCursor},
		)
	})
}

func (f *Scanner) executeQuery(postFunc func() ([]byte, error)) bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	body, err := postFunc()

	if err != nil {
		f.err = errors.Wrapf(err, "can't fecth folder page %s", f.path)
		return false
	}

	var response internal.ListFolderResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		f.err = errors.Wrap(err, "can't parse folder page entries")
		return false
	}

	f.hasNextPage = &response.HasMore
	f.nextCursor = response.Cursor
	f.index = 0
	f.buffer = make([]Action, len(response.Entries))
	for i, entry := range response.Entries {
		actionType := ActionTypeCreate

		if entry.Tag == "deleted" {
			actionType = ActionTypeDelete
		}

		file := fileFromAPI(&entry)
		file.RelativePath = relativePath(f.path, file.RemotePath)

		f.buffer[i] = Action{Type: actionType, File: *file}
	}

	return len(f.buffer) > 0
}

func (f *Scanner) waitForUpdate() error {
	timeout := 30 // seconds
	f.logger.Debugf("wait for new updates (timeout: %d seconds)", timeout)

	body, err := internal.UnuathenticatedPOSTWithBody(
		"https://notify.dropboxapi.com/2/files/list_folder/longpoll",
		map[string]interface{}{"cursor": f.nextCursor, "timeout": timeout},
	)

	if err != nil {
		return errors.Wrap(err, "failure while waiting for new updates")
	}

	var response internal.LongPollResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.Wrap(err, "can't parse polling response")
	}

	if !response.NewFilesAvailable {
		f.logger.Debugf("no new files available")

		return f.waitForUpdate()
	}

	f.logger.Debugf("new files available")

	return nil
}

func relativePath(base string, path string) string {
	return strings.TrimPrefix(path, base)
}
