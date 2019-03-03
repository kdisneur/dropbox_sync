package version

import (
	"fmt"
	"runtime"
)

var (
	buildDate        = "1970-01-01T00:00:00Z"
	gitBranch        = "unknown"
	gitCommit        = "unknown"
	gitState         = "unknown"
	version          = "0.0.0"
	prereleaseNumber = "0"
	isRelease        = ""
	buildNumber      = ""
)

// Info represents all the information about the command line, from the commit
// to the destination platform
type Info struct {
	BuildDate string
	Compiler  string
	GitBranch string
	GitCommit string
	GitState  string
	GoVersion string
	Platform  string
	Version   string
}

// GetInfo returns the current versoin of the command line
// The informations come from ldflags. You can look at the Makefile
// for more information
func GetInfo() Info {
	fullVersion := version

	if isRelease == "" {
		fullVersion += fmt.Sprintf("-alpha.%s", prereleaseNumber)

		if buildNumber != "" {
			fullVersion += fmt.Sprintf("+%s", buildNumber)
		}
	}

	return Info{
		BuildDate: buildDate,
		Compiler:  runtime.Compiler,
		GitBranch: gitBranch,
		GitCommit: gitCommit,
		GitState:  gitState,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Version:   fullVersion,
	}
}
