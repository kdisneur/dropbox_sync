package local

import (
	"github.com/kdisneur/dropbox_sync/pkg/local/internal"
)

const (
	// ActionTypeCreate represents a file creation / update
	ActionTypeCreate internal.ActionType = "create"
	// ActionTypeDelete represents a file deletion
	ActionTypeDelete internal.ActionType = "delete"
)

// Action represents an action and the associate file happening on the
// local filesystem
type Action struct {
	Type internal.ActionType
	File File
}
