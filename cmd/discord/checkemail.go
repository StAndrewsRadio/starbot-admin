package discord

import (
	"fmt"
	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type cmdCheckEmail struct {
	*CommandManager
}

func (cmdCheckEmail) name() string {
	return "checkemail"
}

func (cmdCheckEmail) description() string {
	return "checks if an email is linked to a user, returning the user it is linked to"
}

func (cmdCheckEmail) syntax() string {
	return "<email address>"
}

func (cmd cmdCheckEmail) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	// perm check
	if utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleModerator)) {
		email := strings.Fields(message.Content)[1]

		// check if email is valid
		if !cmd.IsValidEmail(email) {
			_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgCmdEmailNotInList))
			return err
		}

		registered, user, err := cmd.IsEmailRegistered(email)
		if err != nil {
			return err
		}

		// check if registered
		if !registered {
			_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgCmdEmailNotLinked))
			return err
		}

		// send the user
		_, err = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         fmt.Sprintf(cmd.GetString(cfg.MsgCmdEmailLinked), user),
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		})
		return err
	}

	return nil
}
