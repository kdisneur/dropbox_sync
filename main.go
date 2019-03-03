package main

import (
	"os"

	"github.com/kdisneur/dropbox_sync/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	var debuggingFlag bool
	pflag.BoolVar(&debuggingFlag, "debug", false, "enable debug logging")

	var helpFlag bool
	pflag.BoolVarP(&helpFlag, "help", "h", false, "show the current message")

	var versionFlag bool
	pflag.BoolVarP(&versionFlag, "version", "v", false, "show version number")

	pflag.Parse()

	setupLogger(debuggingFlag)

	if helpFlag {
		pflag.Usage()
		os.Exit(0)
	}

	if versionFlag {
		cmd.Version{}.Run()
		os.Exit(0)
	}

	cmd.Synchronize{}.Run()
}

func setupLogger(withDebugging bool) {
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true, FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

	if withDebugging {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
