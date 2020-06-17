package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/StAndrewsRadio/starbot-admin/jobs"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
)

type cmdUnregister struct {
	*CommandManager
}

func (cmdUnregister) name() string {
	return "unregister"
}

func (cmdUnregister) description() string {
	return "removes a show from the database"
}

func (cmdUnregister) syntax() string {
	return "<day> <hour>"
}

func (cmd cmdUnregister) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	// permission check
	if utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleModerator)) {
		args := strings.Split(message.Content, " ")

		// syntax check
		if len(args) != 3 {
			_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgSyntaxError)+
				cmd.syntax())
			if err != nil {
				return err
			}
		} else {
			day, hour := args[1], args[2]
			_, err := time.Parse(db.TimeFormat, day+" "+hour)

			logrus.WithField("day", day).WithField("hour", hour).Debug("Unregistering show...")

			// check they have put a correct time syntax
			if err != nil {
				_, secondErr := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgInvalidTime))
				if secondErr != nil {
					return secondErr
				}

				return nil
			}

			// check the results of deleting the show
			err = cmd.DeleteShow(day, hour)
			if err == buntdb.ErrNotFound {
				_, secondErr := session.ChannelMessageSend(message.ChannelID,
					cmd.GetString(cfg.MsgCmdUnregisterNotFound))
				if secondErr != nil {
					return secondErr
				}
			} else if err != nil {
				return err
			} else {
				go jobs.UpdateShowsEmbed(session, cmd.Database, cmd.Config)

				_, secondErr := session.ChannelMessageSend(message.ChannelID,
					fmt.Sprintf(cmd.GetString(cfg.MsgCmdUnregisterDeleted), day, hour))
				if secondErr != nil {
					return err
				}
			}
		}
	}

	return nil
}
