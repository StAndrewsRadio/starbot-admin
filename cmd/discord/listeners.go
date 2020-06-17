package discord

import (
	"strings"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
)

type cmdListeners struct {
	*CommandManager
}

func (cmdListeners) name() string {
	return "listeners"
}

func (cmdListeners) description() string {
	return "displays the current listener count"
}

func (cmdListeners) syntax() string {
	return ""
}

func (cmd cmdListeners) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	if utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleOnAir)) ||
		utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleModerator)) {

		listeners, err := utils.ReadFromUrl(cmd.GetString(cfg.MiscCurrentListenersUrl))
		if err != nil {
			return err
		}

		listeners = strings.TrimSpace(listeners)

		if listeners == "1" {
			if _, err := session.ChannelMessageSend(message.ChannelID, "There is currently 1 listener."); err != nil {
				return err
			}
		} else {
			if _, err := session.ChannelMessageSend(message.ChannelID, "There are currently "+listeners+" listeners."); err != nil {
				return err
			}
		}
	}

	return nil
}
