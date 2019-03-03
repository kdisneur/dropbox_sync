package local

import (
	"github.com/kdisneur/dropbox_sync/pkg/local/internal"
	"os"
)

const (
	// FileTypeFolder is a directory
	FileTypeFolder internal.FileType = "folder"
	// FileTypeFile is a regular file
	FileTypeFile internal.FileType = "file"
)

// File represents a file or folder on local system
type File struct {
	Path         string
	RelativePath string
	Type         internal.FileType
}

func fileFromEvent(eventName string) File {
	fileType := fileTypeFromPath(eventName)

	return File{
		Type:         fileType,
		Path:         eventName,
		RelativePath: eventName,
	}
}

func fileTypeFromPath(path string) internal.FileType {
	info, err := os.Stat(path)
	if err != nil {
		return FileTypeFile
	}

	if info.IsDir() {
		return FileTypeFolder
	}

	return FileTypeFile
}
