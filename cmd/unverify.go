package cmd

import (
	"strings"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/tidwall/buntdb"
)

type cmdUnverify struct {
	*CommandManager
}

func (cmdUnverify) name() string {
	return "unverify"
}

func (cmdUnverify) description() string {
	return "removes an email address that has been taken by a user and also removes their valid role"
}

func (cmdUnverify) syntax() string {
	return "<email>"
}

func (cmd cmdUnverify) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	if utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleModerator)) {
		args := strings.Fields(message.Content)

		if len(args) != 2 {
			_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgSyntaxError)+
				cmd.syntax())
			if err != nil {
				return err
			}
		} else {
			userID, err := cmd.InvalidateEmail(args[1])
			if err == buntdb.ErrNotFound {
				_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.VerificationEmailNotFound))
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}

			err = session.GuildMemberRoleRemove(message.GuildID, userID, cmd.GetString(cfg.RoleVerified))
			if err != nil {
				return err
			}

			_, err = session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.VerificationUserUnverified))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
