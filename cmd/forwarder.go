package cmd

import (
	"strings"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Represents a single command.
type Command interface {
	name() string
	description() string
	syntax() string
	handler(session *discordgo.Session, message *discordgo.MessageCreate) error
}

// Manager for all commands.
type CommandManager struct {
	*cfg.Config
	*db.Database

	Prefix       string
	PrefixLength int
	Commands     map[string]Command
}

// Makes a new command manager, filling in all available commands.
func New(config *cfg.Config, database *db.Database) *CommandManager {
	commandMap := make(map[string]Command)

	// construct the manager and return
	mgr := &CommandManager{
		Config:   config,
		Database: database,
		Prefix:   config.GetString(cfg.BotPrefix),
		Commands: commandMap,
	}

	// fill in commands
	commandMap["help"] = cmdHelp{mgr}
	commandMap["register"] = cmdRegister{mgr}
	commandMap["show"] = cmdShow{mgr}
	commandMap["invite"] = cmdInvite{mgr}
	commandMap["uninvite"] = cmdUninvite{mgr}
	commandMap["unregister"] = cmdUnregister{mgr}

	logrus.WithField("cmds", len(commandMap)).Debug("New command manager created!")

	mgr.PrefixLength = len(mgr.Prefix)

	return mgr
}

// Event handler for the message create event that forwards commands to the correct command handler.
func (manager *CommandManager) CommandForwarder(session *discordgo.Session, message *discordgo.MessageCreate) {
	// ignore messages from the bot
	if message.Author.ID == session.State.User.ID {
		return
	}

	// check prefix
	if strings.HasPrefix(message.Content, manager.Prefix) {
		spaceIndex := strings.Index(message.Content, " ")
		if spaceIndex < 0 {
			spaceIndex = len(message.Content)
		}

		// get the command name, try and retrieve an executor, then execute if we can!
		cmdString := utils.Substring(message.Content, manager.PrefixLength, spaceIndex)
		cmdExecutor, ok := manager.Commands[cmdString]
		if ok {
			logrus.WithField("cmdString", cmdString).WithField("cmd", cmdExecutor.name()).
				Debug("Executing command...")

			err, str := cmdExecutor.handler(session, message), "There was an error whilst executing that command!"

			if err != nil {
				_, secondErr := session.ChannelMessageSend(message.ChannelID, str+"\n"+err.Error())
				if secondErr != nil {
					logrus.WithError(secondErr).Error("There was an error whilst sending a message to the server!")
				}

				logrus.WithError(err).Error(str)
			}
		}
	}
}
