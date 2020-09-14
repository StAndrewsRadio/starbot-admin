package jobs

import (
	"sync"
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
)

var (
	onAirRole    string
	newHosts     []string
	handlerAdded bool
	wg           sync.WaitGroup

	swapLogger = logrus.WithField("event", "swapShows")
)

// Kicks out the current show hosts and moves the new key show host in, returns true if there is another show next.
func SwapJob(database *db.Database, session *discordgo.Session, config *cfg.Config) (error, bool) {
	// check the swap handler has been added
	if !handlerAdded {
		logrus.Debug("Registering swap job handlers for chunk received...")

		onAirRole = config.GetString(cfg.RoleOnAir)
		session.AddHandler(swapMemberChunkReceived)
		handlerAdded = true
	}

	// add ten minutes to the current time just in case we're straddling slightly behind an hour
	currentTime := time.Now().Add(10 * time.Minute)

	newShow, err := database.GetShow(currentTime.Format(db.DayFormat), currentTime.Format(db.HourFormat))
	if err != nil && err != buntdb.ErrNotFound {
		swapLogger.WithError(err).Error("An error occurred whilst retrieving the next show from the database.")
		return err, false
	}

	swapLogger.WithField("show", newShow).WithField("time", time.Now()).Debug("Running event...")

	// set or clear the new key host
	if err == buntdb.ErrNotFound {
		newHosts = nil
	} else {
		//noinspection GoNilness (ide being dumb, err must be nil at this point)
		newHosts = newShow.Hosts
	}

	// wait for the chunks to all come through
	wg.Add(1)

	// request the members chunk for our guild
	err = session.RequestGuildMembers(config.GetString(cfg.GeneralGuild), "", 0, false)
	if err != nil {
		swapLogger.WithError(err).Error("An error occurred whilst requesting a list of guild members.")
		return err, false
	}

	logrus.Debug("Waiting for chunk responses...")
	wg.Wait()
	return nil, newHosts != nil
}

// Iterates through all members, removing the role from those that do not match the new key host and adding the role
// to those who do.
func swapMemberChunkReceived(session *discordgo.Session, chunk *discordgo.GuildMembersChunk) {
	logrus.WithField("chcnt", chunk.ChunkCount).WithField("chind", chunk.ChunkIndex).
		Debug("Guild members chunk received!")

	for _, member := range chunk.Members {
		if utils.StringSliceContains(newHosts, member.User.ID) {
			// add the on air role
			err := session.GuildMemberRoleAdd(chunk.GuildID, member.User.ID, onAirRole)
			if err != nil {
				swapLogger.WithError(err).Error("An error occurred whilst adding the on air role to a user.")
			}
		} else if utils.StringSliceContains(member.Roles, onAirRole) {
			// remove the on air role
			err := session.GuildMemberRoleRemove(chunk.GuildID, member.User.ID, onAirRole)
			if err != nil {
				swapLogger.WithError(err).Error("An error occurred whilst removing the on air role from a user.")
			}
		}
	}

	// check chunk index
	if chunk.ChunkIndex+1 >= chunk.ChunkCount {
		logrus.Debug("Final chunk received!")
		wg.Done()
	}
}
