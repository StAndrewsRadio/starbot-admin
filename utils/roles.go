package utils

import (
	"github.com/bwmarrin/discordgo"
)

// Checks if a sender of a message is in a given role.
func IsSenderInRole(session *discordgo.Session, message *discordgo.MessageCreate, role string) bool {
	return IsUserInRole(session, message.GuildID, message.Author.ID, role)
}

// Checks if a user is in a specific role.
func IsUserInRole(session *discordgo.Session, guildID, userID, role string) bool {
	member, err := session.GuildMember(guildID, userID)
	if err != nil {
		return false
	}

	return StringSliceContains(member.Roles, role)
}
