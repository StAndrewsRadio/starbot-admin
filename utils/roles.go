package utils

import (
	"github.com/bwmarrin/discordgo"
)

// Checks if a sender of a message is in a given role.
func IsSenderInRole(session *discordgo.Session, message *discordgo.MessageCreate, role string) bool {
	member, err := session.GuildMember(message.GuildID, message.Author.ID)
	if err != nil {
		return false
	}

	return StringSliceContains(member.Roles, role)
}
