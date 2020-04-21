package jobs

import (
	"fmt"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/StAndrewsRadio/starbot-admin/db"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
)

var (
	days  = [...]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	hours = [...]string{"12", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}
)

func CreateShowsEmbed(channelID string) (string, string, bool, error) {
	embed, err := session.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Title:       "Show Schedule",
		Description: "Coming soon...",
		Color:       0xfdde58,
	})
	if err != nil {
		logrus.WithError(err).Error("An error occurred whilst creating the shows embed message.")
		return "", "", false, nil
	}

	previousChannel, previousMessage, replaced, err := database.SetShowsEmbed(channelID, embed.ID)

	// update the embed in a new thread if there was no db error
	if err == nil {
		go UpdateShowsEmbed()
	} else {
		logrus.WithError(err).Debug("Some sort of error happened during the creation of the embed")
	}

	return previousChannel, previousMessage, replaced, err
}

func UpdateShowsEmbed() {
	logrus.Debug("Updating shows embed...")

	// get the message id
	channelID, messageID, err := database.GetShowsEmbed()
	if err != nil {
		if err == buntdb.ErrNotFound {
			logrus.Warn("The update shows job was called but no message ID is stored in the database.")
		} else {
			logrus.WithError(err).Error("An error occurred whilst getting the message ID during the update" +
				" shows embed job.")
		}

		return
	}

	// create the embed
	embed := &discordgo.MessageEmbed{
		Title: "Show Schedule",
		Description: fmt.Sprintf("For more information about a show you can use the show command. For "+
			"example, to find out more info about the Monday 8PM show, type `%sshow Monday 8PM`.",
			config.GetString(cfg.BotPrefix)),
		Color:  0xfdde58,
		Fields: []*discordgo.MessageEmbedField{},
	}

	tx, err := database.GetRawDatabase().Begin(false)
	if err != nil {
		logrus.WithError(err).Error("An error occurred whilst opening a transaction during the update" +
			" shows embed job.")
		return
	}

	// fill it in
	for _, day := range days {
		field := &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("**%s**", day),
			Value:  "",
			Inline: true,
		}

		// am
		for _, hour := range hours {
			show := &db.Show{
				Day:  day,
				Hour: hour + "AM",
			}

			err := db.FillShowFromTransaction(show, tx)
			if err != nil {
				if err != buntdb.ErrNotFound {
					logrus.WithError(err).Error("An error occurred whilst loading a show for the update embed " +
						"job.")
				}

				field.Value += fmt.Sprintf("**%s**:\n", show.Hour)
			} else {
				field.Value += fmt.Sprintf("**%s**: *%s*\n", show.Hour, show.Name)
			}
		}

		// pm
		for _, hour := range hours {
			show := &db.Show{
				Day:  day,
				Hour: hour + "PM",
			}

			err := db.FillShowFromTransaction(show, tx)
			if err != nil {
				if err != buntdb.ErrNotFound {
					logrus.WithError(err).Error("An error occurred whilst loading a show for the update embed " +
						"job.")
				}

				field.Value += fmt.Sprintf("**%s**:\n", show.Hour)
			} else {
				field.Value += fmt.Sprintf("**%s**: *%s*.\n", show.Hour, show.Name)
			}
		}

		field.Value += string('\u200B')
		embed.Fields = append(embed.Fields, field)
	}

	// close the transaction
	err = tx.Rollback()
	if err != nil {
		logrus.WithError(err).Error("An error occurred whilst closing the transaction during the update" +
			" shows embed.")
	}

	// edit the message
	_, err = session.ChannelMessageEditEmbed(channelID, messageID, embed)
	if err != nil {
		logrus.WithError(err).Error("An error occurred whilst editing the update shows embed.")
	}
}
