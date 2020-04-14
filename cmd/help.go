package cmd

import (
	"github.com/bwmarrin/discordgo"
)

type cmdHelp struct {
	*CommandManager
}

func (cmdHelp) name() string {
	return "help"
}

func (cmdHelp) description() string {
	return "displays help information"
}

func (cmdHelp) syntax() string {
	return ""
}

func (command cmdHelp) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	help := ""

	for _, value := range command.Commands {
		help += command.Prefix + value.name()

		if value.syntax() != "" {
			help += " " + value.syntax()
		}

		help += ": " + value.description() + "\n"
	}

	_, err := session.ChannelMessageSend(message.ChannelID, help)
	return err
}

