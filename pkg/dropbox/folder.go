package dropbox

import (
	"github.com/kdisneur/dropbox_sync/pkg/dropbox/internal"
)

// FolderCreate creates a folder on Dropbox if not present on Dropbox
func FolderCreate(client Client, path string) error {
	_, err := internal.POSTWithBody(
		"https://api.dropboxapi.com/2/files/create_folder_v2",
		client.token,
		map[string]interface{}{"path": path, "autorename": false},
	)

	return err
}
