package vars

import (
	"math/rand"
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type StarbotOption func(*Starbot)

func WithConfig(path string) StarbotOption {
	return func(starbot *Starbot) {
		starbot.Config, starbot.err = cfg.New(path)

		// check if debug is enabled
		if starbot.err == nil && starbot.Config.GetBool(cfg.GeneralDebug) {
			logrus.SetLevel(logrus.DebugLevel)
		}
	}
}

func WithDatabase(path string) StarbotOption {
	return func(starbot *Starbot) {
		starbot.Database, starbot.err = db.Open(path)

		// check err before final setting as we can try loading from config
		if starbot.err != nil {
			// only log it as an error if the path was actually set to something
			if path != "" {
				logrus.WithError(starbot.err).Error("The database could not be loaded from the command line argument!")
			}

			starbot.err = nil
			starbot.getDatabaseFromConfig = true
		}
	}
}

func WithBotSession(args ...interface{}) StarbotOption {
	return func(starbot *Starbot) {
		if len(args) == 0 {
			starbot.getBotSessionFromConfig = true
		} else {
			starbot.BotSession, starbot.err = discordgo.New(args...)
		}
	}
}

func WithUserSession(args ...interface{}) StarbotOption {
	return func(starbot *Starbot) {
		if len(args) == 0 {
			starbot.getUserSessionFromConfig = true
		} else {
			starbot.UserSession, starbot.err = discordgo.New(args...)
		}
	}
}

func WithEmailer() StarbotOption {
	return func(starbot *Starbot) {
		starbot.withEmailer = true
	}
}

func WithCommander() StarbotOption {
	return func(starbot *Starbot) {
		starbot.withCommander = true
	}
}

func WithSeededRandom() StarbotOption {
	return func(starbot *Starbot) {
		rand.Seed(time.Now().UnixNano())
	}
}

func WithTriggers() StarbotOption {
	return func(starbot *Starbot) {
		starbot.withTriggers = true
	}
}

func AwaitUserExit() StarbotOption {
	return func(starbot *Starbot) {
		starbot.awaitUserExit = true
	}
}
