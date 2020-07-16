package discord

import (
	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/jobs"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type cmdAutoplay struct {
	*CommandManager
}

func (cmdAutoplay) name() string {
	return "autoplay"
}

func (cmdAutoplay) description() string {
	return "manually starts autoplay"
}

func (cmdAutoplay) syntax() string {
	return ""
}

func (cmd cmdAutoplay) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	if utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleModerator)) {
		_, err := session.ChannelMessageSend(message.ChannelID, "Starting autoplay...")
		if err != nil {
			logrus.WithField("cmd", "autoplay").WithError(err).Error("An error occurred whilst " +
				"sending a message.")
		}

		// start the autoplay
		jobs.StartAutoplay(session, cmd.UserSession, cmd.Config, true, false)
	}

	return nil
}
