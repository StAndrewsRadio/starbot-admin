package vars

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/cmd/discord"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/StAndrewsRadio/starbot-admin/triggers"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type Starbot struct {
	Config      *cfg.Config
	Database    *db.Database
	Commander   *discord.CommandManager
	BotSession  *discordgo.Session
	UserSession *discordgo.Session
	Emailer     *utils.Emailer

	err error
	withEmailer, withCommander, withTriggers, getDatabaseFromConfig, awaitUserExit, withJobs,
	getUserSessionFromConfig, getBotSessionFromConfig bool
}

// Creates a new instance of the Starbot.
func InitialiseStarbot(opts ...StarbotOption) (*Starbot, error) {
	// create the default starbot
	starbot := &Starbot{
		withEmailer:           false,
		withCommander:         false,
		withTriggers:          false,
		getDatabaseFromConfig: false,
		awaitUserExit:         false,
		withJobs:              false,
	}

	// iterate through all options and set them
	for _, opt := range opts {
		opt(starbot)

		if starbot.err != nil {
			return nil, starbot.err
		}
	}

	// set up dependencies
	if err := checkDependencies(starbot); err != nil {
		return nil, err
	}

	// try and open the sessions if they aren't nil
	if starbot.UserSession != nil {
		starbot.UserSession.ShouldReconnectOnError = true
		starbot.UserSession.StateEnabled = false
		starbot.UserSession.Identify.Intents = nil
		logrus.Info("Opening user session...")
		err := starbot.UserSession.Open()
		if err != nil {
			return nil, err
		}
	}

	if starbot.BotSession != nil {
		starbot.BotSession.ShouldReconnectOnError = true
		starbot.BotSession.StateEnabled = true
		starbot.BotSession.Identify.Intents = nil
		logrus.Info("Opening bot session...")
		err := starbot.BotSession.Open()
		if err != nil {
			return nil, err
		}
	}

	// initialise the trigger system
	if starbot.withTriggers {
		go triggers.SetupTriggers(starbot.BotSession, starbot.UserSession, starbot.Config)
	}

	// clean exit
	if starbot.awaitUserExit {
		awaitCleanExit(starbot)
	}

	return starbot, nil
}

func checkDependencies(starbot *Starbot) error {
	if starbot.getUserSessionFromConfig {
		if starbot.Config == nil {
			return errors.New("the user session could not loaded from the config file as the config file was " +
				"not set")
		} else if starbot.UserSession, starbot.err =
			discordgo.New(starbot.Config.GetString(cfg.UserEmail), starbot.Config.GetString(cfg.UserPassword),
				starbot.Config.GetString(cfg.UserToken)); starbot.err != nil {
			return starbot.err
		}
	}

	if starbot.getBotSessionFromConfig {
		if starbot.Config == nil {
			return errors.New("the bot session could not loaded from the config file as the config file was " +
				"not set")
		} else if starbot.BotSession, starbot.err =
			discordgo.New("Bot " + starbot.Config.GetString(cfg.BotToken)); starbot.err != nil {
			return starbot.err
		}
	}

	if starbot.getDatabaseFromConfig {
		if starbot.Config == nil {
			return errors.New("the database could not loaded from the config file as the config file was not set")
		} else if starbot.Database, starbot.err = db.Open(starbot.Config.GetString(cfg.DbFile)); starbot.err != nil {
			return starbot.err
		}
	}

	if starbot.withEmailer {
		if starbot.Config == nil {
			return errors.New("emailer depends on config being loaded")
		} else if starbot.Emailer, starbot.err = utils.NewEmailer(starbot.Config); starbot.err != nil {
			return starbot.err
		}
	}

	if starbot.withCommander {
		if starbot.Config == nil || starbot.Database == nil || starbot.Emailer == nil || starbot.BotSession == nil ||
			starbot.UserSession == nil {

			logrus.WithField("cfg", starbot.Config == nil).WithField("db", starbot.Database == nil).
				WithField("em", starbot.Emailer == nil).WithField("bs", starbot.BotSession == nil).
				WithField("us", starbot.UserSession == nil).Info("Dependency information.")
			return errors.New("commander depends on config, database, emailer and bot session being loaded")
		} else {
			starbot.Commander = discord.New(starbot.Config, starbot.Database, starbot.Emailer, starbot.UserSession)
			starbot.BotSession.AddHandler(starbot.Commander.CommandForwarder)
		}
	}

	if starbot.withTriggers {
		if starbot.Config == nil || starbot.BotSession == nil || starbot.UserSession == nil {
			logrus.WithField("cfg", starbot.Config == nil).WithField("us", starbot.UserSession == nil).
				WithField("bs", starbot.Emailer == nil).Info("Dependency information.")
			return errors.New("the triggers system could not be loaded as the config and Discord sessions " +
				"were not loaded")
		}
	}

	if starbot.withJobs {
		if starbot.Config == nil || starbot.Database == nil || starbot.UserSession == nil || starbot.BotSession == nil {
			logrus.WithField("cfg", starbot.Config == nil).WithField("us", starbot.UserSession == nil).
				WithField("bs", starbot.Emailer == nil).WithField("db", starbot.Database == nil).
				Info("Dependency information.")
			return errors.New("the jobs system could not be loaded without with config, database and Discord " +
				"sessions being loaded")
		}
	}

	return nil
}

func awaitCleanExit(starbot *Starbot) {
	// wait to be killed or terminated before cleanly closing everything
	logrus.Info("The program is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println()

	if starbot.withTriggers {
		logrus.Info("Closing triggers...")
		triggers.Close()
	}

	if starbot.Database != nil {
		logrus.Info("Closing database...")
		if err := starbot.Database.Close(); err != nil {
			logrus.WithError(err).Error("There was an error whilst closing the database!")
		}
	}

	if starbot.BotSession != nil {
		logrus.Info("Closing bot session...")
		if err := starbot.BotSession.Close(); err != nil {
			logrus.WithError(err).Error("There was an error whilst closing the bot connection to Discord!")
		}
	}

	if starbot.UserSession != nil {
		if starbot.Config != nil {
			voice := starbot.UserSession.VoiceConnections[starbot.Config.GetString(cfg.GeneralGuild)]
			if voice != nil {
				logrus.Info("Closing user voice sessions...")
				voice.Close()

				if err := voice.Disconnect(); err != nil {
					logrus.WithError(err).Error("An error occurred whilst closing the current user voice session.")
				}
			}
		}

		logrus.Info("Closing user session...")
		if err := starbot.UserSession.Close(); err != nil {
			logrus.WithError(err).Error("There was an error whilst closing the user connection to Discord!")
		}
	}

}
