package local

import (
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

// Scanner represents a list of actions
type Scanner struct {
	actionEvents  chan Action
	currentAction *Action
	err           error
	errEvents     chan error
	logger        *logrus.Entry
	path          string
	watcher       *fsnotify.Watcher
}

// NewScanner creates a new folder scanner
func NewScanner(logger *logrus.Entry, path string) *Scanner {
	scanner := &Scanner{
		logger:       logger,
		errEvents:    make(chan error),
		actionEvents: make(chan Action),
		path:         path,
	}

	watcher, err := fsnotify.NewWatcher()
	scanner.watcher = watcher
	if err != nil {
		scanner.err = err
		return scanner
	}

	go scanner.listenEvents()

	err = watcher.Add(path)
	if err != nil {
		scanner.err = err
	}

	return scanner
}

// Next replace the `Action` with the following one if it can and return false if it can't
func (s *Scanner) Next() bool {
	if s.err != nil {
		return false
	}

	select {
	case action := <-s.actionEvents:
		s.currentAction = &action
		return true
	case err := <-s.errEvents:
		s.err = err
		return false
	}
}

// NotifyCreation creates a new watcher on the folder
func (s *Scanner) NotifyCreation(relativePath string) {
	absolutePath := path.Join(s.path, relativePath)
	info, err := os.Stat(absolutePath)
	if err != nil {
		return
	}

	if info.IsDir() {
		s.watcher.Add(absolutePath)
	}
}

// NotifyDeletion removes the watcher of the folder
func (s *Scanner) NotifyDeletion(relativePath string) {
	absolutePath := path.Join(s.path, relativePath)
	s.watcher.Remove(absolutePath)
}

// Entry returns the current scanner content
func (s *Scanner) Entry() *Action {
	if s.Err() != nil {
		return nil
	}

	return s.currentAction
}

// Err returns the scanner error if one exists
func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) listenEvents() {
	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				s.errEvents <- errors.New("invalid event received")
			}

			file := fileFromEvent(event.Name)
			file.RelativePath = relativePath(s.path, file.Path)

			action := Action{Type: ActionTypeCreate, File: file}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				action.Type = ActionTypeDelete
			}

			s.actionEvents <- action
		case err, ok := <-s.watcher.Errors:
			if ok {
				s.errEvents <- err
			} else {
				s.errEvents <- errors.New("invalid error event received")
			}
		}
	}
}

func relativePath(base string, path string) string {
	return strings.TrimPrefix(path, base)
}
