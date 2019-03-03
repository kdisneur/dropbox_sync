package cmd

import (
	"fmt"

	"github.com/kdisneur/dropbox_sync/pkg/version"
)

// Version prints the command line version
type Version struct{}

// Run prints the command line version
func (v Version) Run() {
	fmt.Printf("%#+v\n", version.GetInfo())
}
