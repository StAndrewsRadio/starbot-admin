package main

import (
	"os"

	"github.com/StAndrewsRadio/starbot-admin/cmd/program"
	"github.com/StAndrewsRadio/starbot-admin/vars"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

var GlobalOptions struct {
	Debug bool `short:"d" long:"debug" description:"Displays debug messages"`
}

func init() {
	for _, arg := range os.Args {
		if arg == "-d" || arg == "--debug" {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("Debug logging enabled!")
		}
	}

	logrus.Debug("Loading parser...")
	vars.Parser = flags.NewParser(&GlobalOptions, flags.Default)

	program.RegisterCommands()
}

func main() {
	if _, err := vars.Parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
