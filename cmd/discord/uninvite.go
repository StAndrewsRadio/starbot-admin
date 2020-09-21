package discord

import (
	"fmt"
	"strings"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
)

type cmdUninvite struct {
	*CommandManager
}

func (cmdUninvite) name() string {
	return "uninvite"
}

func (cmdUninvite) description() string {
	return "removes a user from the studio"
}

func (cmdUninvite) syntax() string {
	return "<@user>"
}

func (cmd cmdUninvite) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	// permission check
	if utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleOnAir)) ||
		utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleModerator)) {
		// syntax check
		args := strings.Fields(message.Content)
		if len(args) != 2 || len(message.Mentions) != 1 {
			_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgSyntaxError)+
				cmd.syntax())
			if err != nil {
				return err
			}
		} else {
			user := message.Mentions[0]

			err := session.GuildMemberRoleRemove(cmd.GetString(cfg.GeneralGuild), user.ID, cmd.GetString(cfg.RoleOnAir))
			if err != nil {
				return err
			}

			_, err = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf(cmd.GetString(cfg.MsgCmdUninvite),
				user.ID))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
