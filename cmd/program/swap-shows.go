package program

import (
	"github.com/StAndrewsRadio/starbot-admin/jobs"
	"github.com/StAndrewsRadio/starbot-admin/vars"
	"github.com/sirupsen/logrus"
)

type SwapShowsArgs struct {
	ConfigFile string `positional-arg-name:"config" description:"The configuration file"`
}

type SwapShowsCommand struct {
	Arguments SwapShowsArgs `positional-args:"true" required:"true"`
	Database  string        `short:"b" long:"database" description:"If set, the database will be read from this file rather than the config file" optional:"true"`
}

var swapShowsCommand SwapShowsCommand

func init() {
	Commands = append(Commands, func() {
		_, err := vars.Parser.AddCommand("swap-shows", "Swaps the current show to the next "+
			"show, triggering autoplay if needed", "Swaps the current show to the next show, triggering "+
			"autoplay if needed", &swapShowsCommand)
		if err != nil {
			logrus.WithError(err).Fatal("An error occurred whilst parsing the swap-shows command!")
		}
	})
}

func (cmd *SwapShowsCommand) Execute(args []string) error {
	starbot, err := vars.InitialiseStarbot(
		vars.WithConfig(cmd.Arguments.ConfigFile),
		vars.WithDatabase(cmd.Database),
		vars.WithBotSession(),
		vars.WithUserSession(),
	)

	if err != nil {
		return err
	}

	logrus.Info("Running swap shows job...")
	err, upcoming := jobs.SwapJob(starbot.Database, starbot.BotSession, starbot.Config)
	if err != nil {
		return err
	}

	if !upcoming {
		logrus.Info("As there's no show next we'll start autoplay...")
		jobs.StartAutoplay(starbot.BotSession, starbot.UserSession, starbot.Config, true)
	}

	logrus.Info("Job executed successfully!")
	return nil
}
