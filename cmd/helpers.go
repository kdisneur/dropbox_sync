package cmd

import (
	"fmt"
	"os"
)

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
