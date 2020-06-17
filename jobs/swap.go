package jobs

import (
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
)

var (
	onAirRole    string
	newKeyHost   string
	handlerAdded bool

	swapLogger = logrus.WithField("event", "swapShows")
)

// Kicks out the current show hosts and moves the new key show host in.
func swapJob(database *db.Database, session *discordgo.Session, config *cfg.Config) {
	// check the swap handler has been added
	if !handlerAdded {
		onAirRole = config.GetString(cfg.RoleOnAir)
		session.AddHandler(swapMemberChunkReceived)
		handlerAdded = true
	}

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
