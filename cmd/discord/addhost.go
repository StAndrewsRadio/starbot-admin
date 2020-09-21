package discord

import (
	"fmt"
	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/tidwall/buntdb"
	"strings"
	"time"
)

type cmdAddHost struct {
	*CommandManager
}

func (cmdAddHost) name() string {
	return "addhost"
}

func (cmdAddHost) description() string {
	return "adds a host to a show"
}

func (cmdAddHost) syntax() string {
	return "<show day> <show time> <new show host>"
}

func (cmd cmdAddHost) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	args := strings.Fields(message.Content)

	// check syntax and correct mentioning
	if len(args) != 4 || len(message.Mentions) != 1 {
		_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgSyntaxError)+cmd.syntax())
		if err != nil {
			return err
		}
	} else {
		user, day, hour := message.Mentions[0].ID, args[1], args[2]

		// check date
		_, err := time.Parse(db.TimeFormat, day+" "+hour)
		if err != nil {
			_, secondErr := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgInvalidTime))
			if secondErr != nil {
				return secondErr
			} else {
				return err
			}
		}

		// check show exists
		show, err := cmd.GetShow(day, hour)
		if err == buntdb.ErrNotFound {
			_, secondErr := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgCmdShowNotFound))
			if secondErr != nil {
				return secondErr
			} else {
				return err
			}
		} else if err != nil {
			return err
		}

		// check that the user is a host for that show or they are a mod
		if !utils.StringSliceContains(show.Hosts, message.Author.ID) &&
			!utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleModerator)) {

			_, secondErr := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgCmdAddHostNotHost))
			if secondErr != nil {
				return secondErr
			} else {
				return err
			}
		}

		// time to add them as a host
		show.Hosts = append(show.Hosts, user)
		_, _, err = cmd.PutShow(show)
		if err != nil {
			return err
		}

		// give them the role
		err = session.GuildMemberRoleAdd(cmd.GetString(cfg.GeneralGuild), user, cmd.GetString(cfg.RoleVerified))
		if err != nil {
			return err
		}

		// let the user know!
		_, err = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:         fmt.Sprintf(cmd.GetString(cfg.MsgCmdAddHostDone), utils.FormatUserList(show.Hosts)),
			AllowedMentions: &discordgo.MessageAllowedMentions{}})
		if err != nil {
			return err
		}
	}

	return nil
}
