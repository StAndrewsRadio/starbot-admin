package cmd

import (
	"fmt"
	"strings"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
)

type cmdInvite struct {
	*CommandManager
}

func (cmdInvite) name() string {
	return "invite"
}

func (cmdInvite) description() string {
	return "invites another user to join the studio"
}

func (cmdInvite) syntax() string {
	return "<@user>"
}

func (cmdInvite cmdInvite) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	// permission check
	if utils.IsSenderInRole(session, message, cmdInvite.GetString(cfg.RoleOnAir)) ||
		utils.IsSenderInRole(session, message, cmdInvite.GetString(cfg.RoleModerator)) {
		// syntax check
		args := strings.Split(message.Content, " ")
		if len(args) != 2 || len(message.Mentions) != 1 {
			_, err := session.ChannelMessageSend(message.ChannelID, cmdInvite.GetString(cfg.MsgSyntaxError)+
				cmdInvite.syntax())
			if err != nil {
				return err
			}
		} else {
			user := message.Mentions[0]

			err := session.GuildMemberRoleAdd(message.GuildID, user.ID, cmdInvite.GetString(cfg.RoleOnAir))
			if err != nil {
				return err
			}

			_, err = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf(cmdInvite.GetString(cfg.MsgCmdInvite),
				user.ID))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
