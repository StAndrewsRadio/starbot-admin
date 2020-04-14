package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/bwmarrin/discordgo"
	"github.com/tidwall/buntdb"
)

type cmdShow struct{
	*CommandManager
}

func (cmdShow) name() string {
	return "show"
}

func (cmdShow) description() string {
	return "gets information about a show"
}

func (cmdShow) syntax() string {
	return "<day> <hour>"
}

func (cmdShow cmdShow) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	args := strings.Split(message.Content, " ")

	// syntax check
	if len(args) != 3 {
		_, err := session.ChannelMessageSend(message.ChannelID, cmdShow.GetString(cfg.MsgSyntaxError) +
			cmdShow.syntax())
		if err != nil {
			return err
		}
	} else {
		_, err := time.Parse(db.TimeFormat, args[1] + " " + args[2])

		// check time validity
		if err != nil {
			_, secondErr := session.ChannelMessageSend(message.ChannelID, cmdShow.GetString(cfg.MsgInvalidTime))
			if secondErr != nil {
				return secondErr
			}
		} else {
			show, err := cmdShow.GetShow(args[1], args[2])
			var text string

			if err == buntdb.ErrNotFound {
				// no show found
				text = cmdShow.GetString(cfg.MsgCmdShowNotFound)
			} else if err != nil {
				// another error occurred
				return err
			} else {
				// we're good to show the info
				text = fmt.Sprintf(cmdShow.GetString(cfg.MsgCmdShowFound), show.Day, show.Hour, show.Name,
					show.KeyHost)
			}

			_, err = session.ChannelMessageSend(message.ChannelID, text)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

