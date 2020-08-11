package main

import (
	"os"

	"github.com/StAndrewsRadio/starbot-admin/cmd/program"
	"github.com/StAndrewsRadio/starbot-admin/vars"
	"github.com/bwmarrin/discordgo"
	"github.com/jessevdk/go-flags"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var GlobalOptions struct {
	Debug   bool   `short:"d" long:"debug" description:"Displays debug messages"`
	LogFile string `short:"l" long:"log-file" description:"Writes output to the file in addition to stdout"`
}

func logWrapper(msgL, caller int, format string, a ...interface{}) {
	logrus.Warnf(format, a)
}

func init() {
	for num, arg := range os.Args {
		if arg == "-d" || arg == "--debug" {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("Debug logging enabled!")
		} else if arg == "-l" || arg == "--log-file" {
			logrus.AddHook(lfshook.NewHook(os.Args[num+1], &logrus.TextFormatter{}))

			logrus.Infof("Logging to \"%s\".", os.Args[num+1])
		}
	}

	discordgo.Logger = logWrapper

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
