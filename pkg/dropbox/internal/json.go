package internal

// AccessTokenResponse represents ths JSON we get back from Dropbox
type AccessTokenResponse struct {
	Token string `json:"access_token"`
}

// ListFolderResponse represents the JSON we get back from Dropbox
// https://www.dropbox.com/developers/documentation/http/documentation#files-list_folder
type ListFolderResponse struct {
	Entries []FileMetadataResponse `json:"entries"`
	Cursor  string                 `json:"cursor"`
	HasMore bool                   `json:"has_more"`
}

// FileMetadataResponse represents a specific JSON entry we get back from Dropbox
// https://www.dropbox.com/developers/documentation/http/documentation#files-list_folder
type FileMetadataResponse struct {
	ID          string `json:"id"`
	Tag         string `json:".tag"`
	Name        string `json:"name"`
	Path        string `json:"path_display"`
	ContentHash string `json:"content_hash"`
}

// LongPollResponse represents the JSON we get back from Dropbox
// https://www.dropbox.com/developers/documentation/http/documentation#files-list_folder-longpoll
type LongPollResponse struct {
	NewFilesAvailable bool `json:"changes"`
}
