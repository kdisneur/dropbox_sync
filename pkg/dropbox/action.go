package dropbox

import (
	"github.com/kdisneur/dropbox_sync/pkg/dropbox/internal"
)

const (
	// ActionTypeCreate represents a new file or folder creation
	ActionTypeCreate internal.ActionType = "create"

	// ActionTypeDelete represents a new file or folder deletion
	ActionTypeDelete internal.ActionType = "delete"
)

// Action represents a Dropbox action
type Action struct {
	Type internal.ActionType
	File File
}
