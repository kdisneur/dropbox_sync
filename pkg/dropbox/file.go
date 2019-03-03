package dropbox

import (
	"encoding/json"
	"github.com/kdisneur/dropbox_sync/pkg/dropbox/internal"
)

const (
	// FileTypeFolder is a directory
	FileTypeFolder internal.FileType = "folder"
	// FileTypeFile is a regular file
	FileTypeFile internal.FileType = "file"
)

// File represents a file or folder on Dropbox
type File struct {
	ID           string
	ContentHash  string
	Name         string
	RelativePath string
	RemotePath   string
	Type         internal.FileType
}

// FileDelete deletes a file if present on Dropbox
func FileDelete(client Client, path string) error {
	_, err := FileMetadata(client, path)
	if err != nil {
		return nil
	}

	_, err = internal.POSTWithBody(
		"https://api.dropboxapi.com/2/files/delete_v2",
		client.token,
		map[string]interface{}{"path": path},
	)

	return err
}

// FileDownload downloads a file from the user's Dropbox
func FileDownload(client Client, path string) ([]byte, error) {
	return internal.POSTWithDataHeaders(
		"https://content.dropboxapi.com/2/files/download",
		client.token,
		map[string]interface{}{"path": path},
	)
}

// FileMetadata fetches file metadata from Dropbox
func FileMetadata(client Client, path string) (*File, error) {
	body, err := internal.POSTWithBody(
		"https://api.dropboxapi.com/2/files/get_metadata",
		client.token,
		map[string]interface{}{"path": path, "include_deleted": false},
	)

	if err != nil {
		return nil, err
	}

	response := &internal.FileMetadataResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	return fileFromAPI(response), nil
}

// FileUpload uploads a file to Dropbox
func FileUpload(client Client, remotePath string, content []byte) error {
	file, err := FileMetadata(client, remotePath)
	if err == nil {
		localHash, errHash := HashFromBytes(content)
		if errHash == nil && localHash == file.ContentHash {
			return nil
		}
	}

	_, err = internal.POSTWithDataHeadersAndBinary(
		"https://content.dropboxapi.com/2/files/upload",
		client.token,
		map[string]interface{}{"path": remotePath, "mode": "add", "autorename": false, "mute": false},
		content,
	)

	return err
}

func fileFromAPI(entry *internal.FileMetadataResponse) *File {
	fileType := FileTypeFile

	if entry.Tag == "folder" {
		fileType = FileTypeFolder
	}

	file := File{
		ID:           entry.ID,
		ContentHash:  entry.ContentHash,
		Name:         entry.Name,
		RelativePath: entry.Path,
		RemotePath:   entry.Path,
		Type:         fileType,
	}

	return &file
}
