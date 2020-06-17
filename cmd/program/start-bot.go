package program

import (
	"github.com/StAndrewsRadio/starbot-admin/vars"
	"github.com/sirupsen/logrus"
)

type StartBotArgs struct {
	ConfigFile string `positional-arg-name:"config" description:"The configuration file"`
}

type StartBotCommand struct {
	Arguments StartBotArgs `positional-args:"true" required:"true"`
	Database  string       `short:"b" long:"database" description:"If set, the database will be read from this file rather than the config file" optional:"true"`
}

var startBotCommand StartBotCommand

func init() {
	Commands = append(Commands, func() {
		_, err := vars.Parser.AddCommand("start-bot", "Starts the Discord bot.",
			"Starts the Discord bot.", &startBotCommand)
		if err != nil {
			logrus.WithError(err).Fatal("An error occurred whilst parsing the start-bot command!")
		}
	})
}

func (cmd *StartBotCommand) Execute(args []string) error {
	_, err := vars.InitialiseStarbot(
		vars.WithSeededRandom(),
		vars.WithConfig(cmd.Arguments.ConfigFile),
		vars.WithEmailer(),
		vars.WithDatabase(cmd.Database),
		vars.WithCommander(),
		vars.WithTriggers(),
		vars.AwaitUserExit(),
		vars.WithJobs(),
		vars.WithBotSession(),
		vars.WithUserSession(),
	)

	return err
}
