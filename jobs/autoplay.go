package jobs

import (
	"fmt"
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	autoplayLogger = logrus.WithField("event", "autoplayJob")
	join           *discordgo.VoiceConnection
	running        = false
)

// Checks if the studio voice channel is empty and plays some music if it is.
func StartAutoplay(session, userSession *discordgo.Session, config *cfg.Config, ignoreUsers, isSlotUp bool) {
	if running {
		autoplayLogger.Warning("Something triggered the job whilst it was already running!")
		return
	}

	running = true

	// log event
	autoplayLogger.WithField("time", time.Now()).Debug("Running event...")

	// get the guild
	guild, err := session.State.Guild(config.GetString(cfg.GeneralGuild))
	if err != nil {
		autoplayLogger.WithError(err).Error("An error occurred whilst getting the guild!")
		running = false
		return
	}

	numInStudio := 0
	studioID, forwarderID := config.GetString(cfg.ChannelStudio), config.GetStrings(cfg.AutoplayIgnoredUsers)
	isUserBotInChannel := false

	for _, voiceState := range guild.VoiceStates {
		if voiceState.ChannelID == studioID {
			// if the state represents a user in the studio that isn't a forwarder, increment the number
			if !utils.StringSliceContains(forwarderID, voiceState.UserID) {
				numInStudio++
			}

			// if the state represents the user bot, then we won't re-join and decrement the num
			if voiceState.UserID == userSession.State.User.ID {
				isUserBotInChannel = true
				numInStudio--
			}
		}
	}

	// if nobody is in the studio, it's autoplayJob time!
	controlRoomID := config.GetString(cfg.ChannelControlRoom)
	if numInStudio == 0 || ignoreUsers {
		// join the studio
		if join == nil && !isUserBotInChannel {
			join, err = userSession.ChannelVoiceJoin(guild.ID, studioID, true, false)
			if err != nil {
				autoplayLogger.WithError(err).Error("An error occurred whilst joining the studio!")
				running = false
				return
			}
		}

		// wait a mo
		time.Sleep(3 * time.Second)

		// announce that we're starting autoplayJob
		var announcementMessage string
		if isSlotUp {
			announcementMessage = config.GetString(cfg.AutoplaySlotUp)
		} else {
			announcementMessage = config.GetString(cfg.AutoplayAnnounce)
		}

		_, err = userSession.ChannelMessageSend(controlRoomID, announcementMessage)
		if err != nil {
			autoplayLogger.WithError(err).Error("An error occurred whilst sending the announcement message.")
		}

		// iterate through all commands
		for _, command := range config.GetStrings(cfg.AutoplayCommands) {
			// wait a mo
			time.Sleep(5 * time.Second)

			// send the command
			_, err := userSession.ChannelMessageSend(controlRoomID, command)
			if err != nil {
				autoplayLogger.WithField("cmd", command).WithError(err).
					Error("An error occurred whilst sending a command.")
			}
		}
	} else {
		// there's people in the studio but someone requested auto play?
		_, err := session.ChannelMessageSendComplex(controlRoomID, &discordgo.MessageSend{
			Content:         fmt.Sprintf(config.GetString(cfg.AutoplayUsersInStudio), config.GetString(cfg.RoleModerator)),
			AllowedMentions: &discordgo.MessageAllowedMentions{}})
		if err != nil {
			autoplayLogger.WithError(err).
				Error("An error occurred whilst sending a message.")
		}
	}

	running = false
}
