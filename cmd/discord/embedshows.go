package discord

import (
	"fmt"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/jobs"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
)

type cmdEmbedShows struct {
	*CommandManager
}

func (cmdEmbedShows) name() string {
	return "embedshows"
}

func (cmdEmbedShows) description() string {
	return "makes or moves the shows information embed to the current channel"
}

func (cmdEmbedShows) syntax() string {
	return ""
}

func (cmd cmdEmbedShows) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	if utils.IsSenderInRole(session, message, cmd.Config.GetString(cfg.RoleModerator)) {
		previousChannel, _, replaced, err := jobs.CreateShowsEmbed(message.ChannelID)
		if err != nil {
			return err
		}

		if replaced {
			_, err = session.ChannelMessageSend(message.ChannelID,
				fmt.Sprintf(cmd.GetString(cfg.MsgCmdShowsEmbedReplaced), previousChannel))
		} else {
			_, err = session.ChannelMessageSend(message.ChannelID,
				fmt.Sprintf(cmd.GetString(cfg.MsgCmdShowsEmbedSet)))
		}

		return err
	}

	return nil
}
