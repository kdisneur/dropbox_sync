package dropbox

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/pkg/errors"
)

// ContentHashBlockSize size of Dropbox hash
const ContentHashBlockSize = 4 * 1024 * 1024

// ContentHasher represents a Dropbox content hasher
// https://www.dropbox.com/developers/reference/content-hash
type ContentHasher struct {
	err             error
	overallChecksum hash.Hash
	partialChecksum [32]byte
}

// HashFromBytes computes a hash on bytes
func HashFromBytes(content []byte) (string, error) {
	reader := bytes.NewReader(content)

	return HashFromReader(reader)
}

// HashFromFile computes a hash from a file path
func HashFromFile(path string) (string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", errors.Wrapf(err, "can't open file '%s'", path)
	}
	defer reader.Close()

	return HashFromReader(reader)
}

// HashFromReader computes a hash from a Reader
func HashFromReader(reader io.Reader) (string, error) {
	c := &ContentHasher{overallChecksum: sha256.New()}

	return c.computeHash(reader)
}

func (c *ContentHasher) computeHash(reader io.Reader) (string, error) {
	for c.nextPartialChecksum(reader) {
		c.overallChecksum.Write(c.partialChecksum[:])
	}

	return fmt.Sprintf("%x", c.overallChecksum.Sum(nil)), c.err
}

func (c *ContentHasher) nextPartialChecksum(reader io.Reader) bool {
	partialFileContent := make([]byte, ContentHashBlockSize)
	bytesRead, err := reader.Read(partialFileContent)
	if err != nil && err != io.EOF {
		c.err = err
		return false
	}

	if err == io.EOF {
		return false
	}

	c.partialChecksum = sha256.Sum256(partialFileContent[:bytesRead])

	return true
}
