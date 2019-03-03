package dropbox

// Client represents an authenticated user
type Client struct {
	token string
}

// NewClient creates a new Dropbox client from a token
func NewClient(token string) Client {
	return Client{token: token}
}
