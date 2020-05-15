package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/cmd"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/StAndrewsRadio/starbot-admin/jobs"
	"github.com/StAndrewsRadio/starbot-admin/triggers"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	config      *cfg.Config
	database    *db.Database
	commander   *cmd.CommandManager
	session     *discordgo.Session
	userSession *discordgo.Session
	emailer     *utils.Emailer
	err         error
)

func main() {
	// seed random
	rand.Seed(time.Now().UnixNano())

	// parse arguments
	args := os.Args

	// check if they want a token
	if args[1] == "get-token" {
		getToken(args)
		return
	}

	// check standard syntax
	if len(args) < 2 {
		logrus.Fatal("Please provide the configuration file to use")
	}

	// read configuration
	config, err = cfg.New(args[1])
	if err != nil {
		logrus.WithError(err).Fatal("There was an error whilst reading the configuration file!")
	}

	if config.GetBool(cfg.GeneralDebug) || utils.StringSliceContains(os.Args[1:], "--debug") {
		logrus.Info("Enabling debug logging...")
		logrus.SetLevel(logrus.DebugLevel)
	}

	// get the emailer
	emailer, err = utils.NewEmailer(config)
	if err != nil {
		logrus.WithError(err).Fatal("There was an error whilst creating the emailer!")
	}

	// open the database
	database, err = db.Open(config.GetString(cfg.DbFile))
	if err != nil {
		logrus.WithError(err).Fatal("There was an error whilst opening the database!")
	}

	// log the bot in
	session, err = discordgo.New("Bot " + config.GetString(cfg.BotToken))
	if err != nil {
		logrus.WithError(err).Fatal("There was an error whilst creating the Discord bot session!")
	}

	// log the user in
	userSession, err = discordgo.New(config.GetString(cfg.UserEmail), config.GetString(cfg.UserPassword),
		config.GetString(cfg.UserToken))
	if err != nil {
		logrus.WithError(err).Fatal("There was an error whilst creating the Discord user session!")
	}

	// get the command manager
	commander = cmd.New(config, database, emailer)

	// register handlers
	session.AddHandler(commander.CommandForwarder)

	// open the bot session
	err = session.Open()
	if err != nil {
		logrus.WithError(err).Fatal("There was an error whilst opening the bot connection to Discord!")
	}

	// open the user session
	err = userSession.Open()
	if err != nil {
		logrus.WithError(err).Fatal("There was an error whilst opening the user connection to Discord!")
	}

	go jobs.ScheduleEvents(config, database, session, userSession)
	go triggers.SetupTriggers(config.GetString(cfg.TriggersAddress), config.GetString(cfg.TriggersPassword),
		session, userSession, config)

	// wait to be killed or terminated before cleanly closing everything
	logrus.Info("The bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	triggers.Close()

	err = database.Close()
	if err != nil {
		logrus.WithError(err).Error("There was an error whilst closing the database!")
	}

	err = session.Close()
	if err != nil {
		logrus.WithError(err).Error("There was an error whilst closing the bot connection to Discord!")
	}

	voice := userSession.VoiceConnections[config.GetString(cfg.GeneralGuild)]
	if voice != nil {
		voice.Close()
		err = voice.Disconnect()
		if err != nil {
			logrus.WithError(err).Error("An error occurred whilst closing the current user voice session.")
		}
	}

	err = userSession.Close()
	if err != nil {
		logrus.WithError(err).Error("There was an error whilst closing the user connection to Discord!")
	}
}

func getToken(args []string) {
	// check arg length
	if len(args) != 4 {
		logrus.Fatal("Please provide an email and password in order to obtain a token.")
	} else {
		session, err := discordgo.New(args[2], args[3])
		if err != nil {
			logrus.WithError(err).Fatal("An error occurred whilst obtaining the token.")
		}

		logrus.WithField("token", session.Token).Info("Token obtained.")
	}
}
