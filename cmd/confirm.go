package cmd

import (
	"strings"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/bwmarrin/discordgo"
)

type cmdConfirm struct {
	*CommandManager
}

func (cmdConfirm) name() string {
	return "confirm"
}

func (cmdConfirm) description() string {
	return "confirms a verification code"
}

func (cmdConfirm) syntax() string {
	return "<code>"
}

func (cmd cmdConfirm) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	args := strings.Fields(message.Content)

	if len(args) != 2 {
		_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgSyntaxError)+
			cmd.syntax())
		if err != nil {
			return err
		}
	} else {
		code := args[1]
		confirmed, err := cmd.CheckVerification(message.Author.ID, code)
		if err != nil {
			return err
		}

		if !confirmed {
			_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.VerificationNotConfirmed))
			if err != nil {
				return err
			}
		} else {
			err = cmd.ValidateUser(message.Author.ID, code)
			if err != nil {
				return err
			}

			_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.VerificationConfirmed))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
