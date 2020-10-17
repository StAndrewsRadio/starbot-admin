package program

import (
	"github.com/StAndrewsRadio/starbot-admin/vars"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type EchoArgs struct {
	ConfigFile string `positional-arg-name:"config" description:"The configuration file"`
	ChannelID  string `positional-arg-name:"channel-id" description:"The ID of the channel to send the message in"`
	Message    string `positional-arg-name:"message" description:"The message to send"`
}

type EchoCommand struct {
	Arguments     EchoArgs `positional-args:"true" required:"true"`
	AllowMentions bool     `short:"m" long:"allow-mentions" description:"If set, the message sent will trigger mentions" optional:"true"`
}

var echoCommand EchoCommand

func init() {
	Commands = append(Commands, func() {
		_, err := vars.Parser.AddCommand("echo", "Echoes a message.",
			"Sends a message to a given channel in the server, optionally triggering mentions.",
			&echoCommand)
		if err != nil {
			logrus.WithError(err).Fatal("An error occurred whilst parsing the echo command!")
		}
	})
}

func (cmd *EchoCommand) Execute(_ []string) error {
	// init
	starbot, err := vars.InitialiseStarbot(
		vars.WithConfig(cmd.Arguments.ConfigFile),
		vars.WithBotSession(),
	)

	// check if any error happened
	if err != nil {
		return err
	}

	// check if they want to allow mentions
	if cmd.AllowMentions {
		_, err = starbot.BotSession.ChannelMessageSend(cmd.Arguments.ChannelID, cmd.Arguments.Message)
		logrus.Debug("Raw message sent.")
	} else {
		_, err = starbot.BotSession.ChannelMessageSendComplex(cmd.Arguments.ChannelID, &discordgo.MessageSend{
			Content:         cmd.Arguments.Message,
			AllowedMentions: &discordgo.MessageAllowedMentions{}})
		logrus.Debug("Complex message sent.")
	}

	if err != nil {
		logrus.WithError(err).Error("There was an error whilst sending the message.")
	} else {
		logrus.Info("Message sent successfully!")
	}

	return nil
}
