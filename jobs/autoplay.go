package jobs

import (
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
)

var (
	autoplayLogger = logrus.WithField("event", "autoplayJob")
	join           *discordgo.VoiceConnection
	running        = false
)

// Schedules the autoplay job.
func autoplayScheduler(quicker bool) {
	// get time to next run
	autoplayTime := time.Now().Truncate(time.Hour).Add(2 * time.Minute).Add(30 * time.Second)
	if autoplayTime.Before(time.Now()) {
		autoplayTime = autoplayTime.Add(time.Hour)
	}

	// set the time to do the job
	var job *gocron.Job
	if quicker {
		job = gocron.Every(1).Minute().From(gocron.NextTick())
	} else {
		job = gocron.Every(1).Hour().From(&autoplayTime)
	}

	// schedule the job
	startTime := job.NextScheduledTime()
	err := job.Do(StartAutoplay, false)
	if err != nil {
		autoplayLogger.WithError(err).Fatal("An error occurred whilst scheduling the autoplayJob job!")
	}

	autoplayLogger.WithField("start", startTime).Debug("Job scheduled successfully.")
}

// Checks if the studio voice channel is empty and plays some music if it is.
func StartAutoplay(ignoreUsers bool) {
	if running {
		autoplayLogger.Warning("Something triggered the job whilst it was already running!")
		return
	}

	running = true

	// log event
	autoplayLogger.WithField("time", time.Now()).Debug("Running event...")

	// get the guild
	guild, err := session.Guild(config.GetString(cfg.GeneralGuild))
	if err != nil {
		autoplayLogger.WithError(err).Error("An error occurred whilst getting the guild!")
		running = false
		return
	}

	numInStudio := 0
	studioID, forwarderID := config.GetString(cfg.ChannelStudio), config.GetStrings(cfg.AutoplayIgnoredUsers)

	// iterate through every voice state only if we care about the users
	if !ignoreUsers {
		for _, voiceState := range guild.VoiceStates {
			// if the state represents a user in the studio that isn't a forwarder, increment the number
			if voiceState.ChannelID == studioID && !utils.StringSliceContains(forwarderID, voiceState.UserID) &&
				voiceState.UserID != userSession.State.User.ID {
				numInStudio++
			}
		}
	}

	// if nobody is in the studio, it's autoplayJob time!
	controlRoomID := config.GetString(cfg.ChannelControlRoom)
	if numInStudio == 0 {
		// join the studio
		if join == nil {
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
		_, err = userSession.ChannelMessageSend(controlRoomID, config.GetString(cfg.AutoplayAnnounce))
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
		_, err := session.ChannelMessageSend(controlRoomID, "Something has requested me to start autoplay "+
			"but there's people here in the studio. Are you sure you're playing music or talking? Message <@&"+
			config.GetString(cfg.RoleSupport)+"> if you need a hand.")
		if err != nil {
			autoplayLogger.WithError(err).
				Error("An error occurred whilst sending a message.")
		}
	}

	running = false
}
