package discord

import (
	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
)

type cmdEcho struct {
	*CommandManager
}

func (cmdEcho) name() string {
	return "echo"
}

func (cmdEcho) description() string {
	return "echoes a message into the current channel"
}

func (cmdEcho) syntax() string {
	return "<message>"
}

func (cmd cmdEcho) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	args := utils.FieldsN(message.Content, 1)

	// role and arg length check
	if utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleModerator)) {
		// check syntax
		if args == nil || len(args) == 0 {
			_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgSyntaxError)+
				cmd.syntax())
			return err
		}

		// get the right echo message
		// check arg length
		if len(args) <= 1 {
			_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgSyntaxError)+
				cmd.syntax())
			if err != nil {
				return err
			}

			return nil
		}

		// send the message!
		_, err := session.ChannelMessageSend(message.ChannelID, args[1])
		if err != nil {
			return err
		}
	}

	return nil
}
