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

type cmdRemoveHost struct {
	*CommandManager
}

func (cmdRemoveHost) name() string {
	return "removehost"
}

func (cmdRemoveHost) description() string {
	return "removes a host from a show"
}

func (cmdRemoveHost) syntax() string {
	return "<show day> <show time> <show host to remove>"
}

func (cmd cmdRemoveHost) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
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

		// time to remove them as a host
		hostIndex := utils.StringSliceIndexOf(user, show.Hosts)
		if hostIndex != -1 {
			// swap the last index to the found index and knock one off the end of the slice
			show.Hosts[hostIndex] = show.Hosts[len(show.Hosts)-1]
			show.Hosts = show.Hosts[:len(show.Hosts)-1]

			// save the result
			_, _, err = cmd.PutShow(show)
			if err != nil {
				return err
			}

			// note we don't remove the role from the user!!! they might have it from elsewhere and we can't easily
			// find this out - this needs to be done manually
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
