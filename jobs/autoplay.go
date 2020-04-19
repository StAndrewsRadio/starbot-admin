package jobs

import (
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
)

var (
	autoplayLogger = logrus.WithField("event", "autoplayJob")
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
	err := job.Do(autoplayJob)
	if err != nil {
		autoplayLogger.WithError(err).Fatal("An error occurred whilst scheduling the autoplayJob job!")
	}

	autoplayLogger.WithField("start", startTime).Debug("Job scheduled successfully.")
}

// Checks if the studio voice channel is empty and plays some music if it is.
func autoplayJob() {
	// log event
	autoplayLogger.WithField("time", time.Now()).Debug("Running event...")

	// get the guild
	guild, err := session.Guild(config.GetString(cfg.GeneralGuild))
	if err != nil {
		autoplayLogger.WithError(err).Error("An error occurred whilst getting the guild!")
		return
	}

	numInStudio := 0
	studioID, forwarderID := config.GetString(cfg.ChannelStudio), config.GetStrings(cfg.AutoplayIgnoredUsers)

	// iterate through every voice state
	for _, voiceState := range guild.VoiceStates {
		// if the state represents a user in the studio that isn't a forwarder, increment the number
		if voiceState.ChannelID == studioID && !utils.StringSliceContains(forwarderID, voiceState.UserID) &&
			voiceState.UserID != userSession.State.User.ID {
			numInStudio++
		}
	}

	// if nobody is in the studio, it's autoplayJob time!
	if numInStudio == 0 {
		controlRoomID := config.GetString(cfg.ChannelControlRoom)

		// join the studio
		join, err := userSession.ChannelVoiceJoin(guild.ID, studioID, true, false)
		if err != nil {
			autoplayLogger.WithError(err).Error("An error occurred whilst joining the studio!")
			return
		}

		// wait a mo
		time.Sleep(5 * time.Second)

		// announce that we're starting autoplayJob
		_, err = userSession.ChannelMessageSend(controlRoomID, config.GetString(cfg.AutoplayAnnounce))
		if err != nil {
			autoplayLogger.WithError(err).Error("An error occurred whilst sending the announcement message.")
		}

		// iterate through all commands
		for _, command := range config.GetStrings(cfg.AutoplayCommands) {
			// send the command
			_, err := userSession.ChannelMessageSend(controlRoomID, command)
			if err != nil {
				autoplayLogger.WithField("cmd", command).WithError(err).
					Error("An error occurred whilst sending a command.")
			}

			// wait one second
			time.Sleep(5 * time.Second)
		}

		// finally, leave the studio
		err = join.Disconnect()
		if err != nil {
			autoplayLogger.WithError(err).Error("An error occurred whilst leaving the studio.")
		}
	}
}
