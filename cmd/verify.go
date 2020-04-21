package cmd

import (
	"fmt"
	"strings"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type cmdVerify struct {
	*CommandManager
}

func (cmdVerify) name() string {
	return "verify"
}

func (cmdVerify) description() string {
	return "sends a verification code to your email address"
}

func (cmdVerify) syntax() string {
	return "<email>"
}

func (cmd cmdVerify) handler(session *discordgo.Session, message *discordgo.MessageCreate) error {
	if !utils.IsSenderInRole(session, message, cmd.GetString(cfg.RoleVerified)) {
		args := strings.Fields(message.Content)

		if len(args) != 2 {
			_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.MsgSyntaxError)+
				cmd.syntax())
			if err != nil {
				return err
			}
		} else {
			email := args[1]

			if !cmd.IsValidEmail(email) {
				_, err := session.ChannelMessageSend(message.ChannelID, cmd.GetString(cfg.VerificationInvalidEmail))
				if err != nil {
					return err
				}
			} else {
				taken, takenBy, err := cmd.IsEmailRegistered(email)
				if err != nil {
					return err
				}

				// don't do anything if it's already taken
				if taken {
					member, err := session.GuildMember(message.GuildID, takenBy)
					if err != nil {
						return err
					}

					_, err = session.ChannelMessageSend(message.ChannelID,
						fmt.Sprintf(cmd.GetString(cfg.VerificationEmailTaken), member.Nick))
					return err
				}

				code := utils.RandomString(8)
				go cmd.SendVerificationEmail(args[1], cmd.GetString(cfg.VerificationEmailSubject), cmd.Prefix, code)

				logrus.WithField("cmd", "verify").Debug("Email goroutine called.")

				err = cmd.StoreVerificationCode(message.Author.ID, email, code)
				if err != nil {
					return err
				}

				logrus.WithField("cmd", "verify").Debug("Verification code stored.")

				_, err = session.ChannelMessageSend(message.ChannelID,
					fmt.Sprintf(cmd.GetString(cfg.VerificationEmailSent), cmd.Prefix))
				if err != nil {
					return err
				}

				logrus.WithField("cmd", "verify").Debug("Response sent.")
			}
		}
	}

	return nil
}
