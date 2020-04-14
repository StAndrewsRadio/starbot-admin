package jobs

import (
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/bwmarrin/discordgo"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
)

var (
	onAirRole  string
	newKeyHost string

	swapLogger = logrus.WithField("event", "swapShows")
)

// Schedules the swap shows job.
func swapScheduler(quicker bool) {
	// get time to next run
	swapperTime := time.Now().Truncate(time.Hour)
	if swapperTime.Before(time.Now()) {
		swapperTime = swapperTime.Add(time.Hour)
	}

	// set the time to do the job
	var job *gocron.Job
	if quicker {
		job = gocron.Every(1).Minute().From(gocron.NextTick())
	} else {
		job = gocron.Every(1).Hour().From(&swapperTime)
	}

	// schedule the job
	startTime := job.NextScheduledTime()
	err := job.Do(swapJob)
	if err != nil {
		autoplayLogger.WithError(err).Fatal("An error occurred whilst scheduling the autoplayJob job!")
	}

	swapLogger.WithField("start", startTime).Debug("Job scheduled successfully.")

	// add the member chunk event handler
	onAirRole = config.GetString(cfg.RoleOnAir)
	session.AddHandler(swapMemberChunkReceived)
}

// Kicks out the current show hosts and moves the new key show host in.
func swapJob() {
	// add ten minutes to the current time just in case we're straddling slightly behind an hour
	currentTime := time.Now().Add(10 * time.Minute)

	newShow, err := database.GetShow(currentTime.Format(db.DayFormat), currentTime.Format(db.HourFormat))
	if err != nil && err != buntdb.ErrNotFound {
		swapLogger.WithError(err).Error("An error occurred whilst retrieving the next show from the database.")
		return
	}

	swapLogger.WithField("show", newShow).WithField("time", time.Now()).Debug("Running event...")

	// set or clear the new key host
	if err == buntdb.ErrNotFound {
		newKeyHost = ""
	} else {
		//noinspection GoNilness (ide being dumb, err must be nil at this point)
		newKeyHost = newShow.KeyHost
	}

	// request the members chunk for our guild
	err = session.RequestGuildMembers(config.GetString(cfg.GeneralGuild), "", 0)
	if err != nil {
		swapLogger.WithError(err).Error("An error occurred whilst requesting a list of guild members.")
		return
	}
}

// Iterates through all members, removing the role from those that do not match the new key host and adding the role
// to those who do.
func swapMemberChunkReceived(session *discordgo.Session, chunk *discordgo.GuildMembersChunk) {
	for _, member := range chunk.Members {
		if member.User.ID == newKeyHost {
			// add the on air role
			err := session.GuildMemberRoleAdd(chunk.GuildID, member.User.ID, onAirRole)
			if err != nil {
				swapLogger.WithError(err).Error("An error occurred whilst adding the on air role to a user.")
			}
		} else {
			// remove the on air role
			err := session.GuildMemberRoleRemove(chunk.GuildID, member.User.ID, onAirRole)
			if err != nil {
				swapLogger.WithError(err).Error("An error occurred whilst removing the on air role from a user.")
			}
		}
	}
}
