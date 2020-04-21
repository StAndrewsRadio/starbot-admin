package db

import "github.com/tidwall/buntdb"

type Database struct {
	db *buntdb.DB
}

const (
	DayFormat  = "Monday"
	HourFormat = "3PM"
	TimeFormat = DayFormat + " " + HourFormat

	ShowPrefix = "show:"
	HostSuffix = ":host"
	NameSuffix = ":name"

	ShowsEmbedChannel = "embed:shows:channel"
	ShowsEmbedMessage = "embed:shows:message"
)

func Open(path string) (*Database, error) {
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, err
}

func (database *Database) Close() error {
	return database.db.Close()
}

func (database *Database) GetRawDatabase() *buntdb.DB {
	return database.db
}

// Sets the shows embed id, returning the old id if it exists.
func (database *Database) SetShowsEmbed(channelID, messageID string) (string, string, bool, error) {
	tx, err := database.db.Begin(true)
	if err != nil {
		return "", "", false, err
	}

	previousChannel, replaced, err := tx.Set(ShowsEmbedChannel, channelID, nil)
	previousMessage, replaced, err := tx.Set(ShowsEmbedMessage, messageID, nil)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return "", "", false, err
		}

		return "", "", false, err
	}

	err = tx.Commit()
	if err != nil {
		return "", "", false, err
	}

	if replaced {
		return previousChannel, previousMessage, true, err
	} else {
		return "", "", false, err
	}
}

// Gets the channel for the shows embed, if it has been set
func (database *Database) GetShowsEmbed() (string, string, error) {
	tx, err := database.db.Begin(false)
	if err != nil {
		return "", "", err
	}

	channelID, err := tx.Get(ShowsEmbedChannel)
	messageID, err := tx.Get(ShowsEmbedMessage)

	err = tx.Rollback()

	return channelID, messageID, err
}
