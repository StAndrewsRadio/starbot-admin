package cmd

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
)

type cmdRegister struct {
	*CommandManager
}

func (cmdRegister) name() string {
	return "register"
}

func (cmdRegister) description() string {
	return "registers a new show, replacing any previous show at that time"
}

func (cmdRegister) syntax() string {
	return "<@host> <show day> <show hour (e.g. 3AM)> <show name>"
}

func (cmdRegister cmdRegister) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	// perm check
	if utils.IsSenderInRole(session, message, cmdRegister.GetString(cfg.RoleModerator)) {
		args := strings.SplitN(message.Content, " ", 5)

		// syntax check
		if len(args) != 5 {
			_, err := session.ChannelMessageSend(message.ChannelID, cmdRegister.GetString(cfg.MsgSyntaxError)+
				cmdRegister.syntax())
			if err != nil {
				return err
			}
		} else {
			// check they've mentioned someone correctly
			if len(message.Mentions) != 1 {
				_, err := session.ChannelMessageSend(message.ChannelID, cmdRegister.GetString(cfg.MsgSyntaxError)+
					cmdRegister.syntax())
				if err != nil {
					return err
				}
			} else {
				user, day, hour, name := message.Mentions[0].ID, args[2], args[3], args[4]
				_, err := time.Parse(db.TimeFormat, day+" "+hour)

				logrus.WithField("day", day).WithField("hour", hour).WithField("name", name).
					WithField("host", user).Debug("Registering show...")

				// check they have put a correct time syntax
				if err != nil {
					_, secondErr := session.ChannelMessageSend(message.ChannelID, cmdRegister.GetString(cfg.MsgInvalidTime))
					if secondErr != nil {
						return secondErr
					}

					return nil
				}

				oldShow, replaced, err := cmdRegister.PutShow(db.Show{
					KeyHost: user,
					Day:     day,
					Hour:    hour,
					Name:    name,
				})
				if err != nil {
					return err
				}

				if replaced {
					_, err = session.ChannelMessageSend(message.ChannelID,
						fmt.Sprintf(cmdRegister.GetString(cfg.MsgCmdRegisterReplaced), user, name, day, hour,
							oldShow.KeyHost, oldShow.Name))
				} else {
					_, err = session.ChannelMessageSend(message.ChannelID,
						fmt.Sprintf(cmdRegister.GetString(cfg.MsgCmdRegisterNewShow), user, name, day, hour))
				}

				go jobs.UpdateShowsEmbed()

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
