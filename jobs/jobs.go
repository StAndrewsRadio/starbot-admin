package jobs

import (
	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/bwmarrin/discordgo"
	"github.com/jasonlvhit/gocron"
)

var (
	config *cfg.Config
	database *db.Database
	session *discordgo.Session
	userSession *discordgo.Session
)

// Schedules all recurring or delayed jobs.
func ScheduleEvents(c *cfg.Config, d *db.Database, s *discordgo.Session, us *discordgo.Session) {
	// set up local variables
	config, database, session, userSession = c, d, s, us
	quicker := config.GetBool(cfg.TestingQuickerJobs)

	// schedule jobs
	swapScheduler(quicker)
	autoplayScheduler(quicker)

	gocron.Start()
}
