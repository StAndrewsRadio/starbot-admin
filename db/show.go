package db

import (
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	"strings"
)

// Represents a single show that is broadcast on STAR.
type Show struct {
	Hosts []string
	Day   string
	Hour  string
	Name  string
}

// Gets a show from the database, given the time it starts.
func (database *Database) GetShow(day, hour string) (Show, error) {
	show := Show{
		Day:  day,
		Hour: hour,
	}

	// load the show from the database
	err := database.db.View(func(tx *buntdb.Tx) error {
		return FillShowFromTransaction(&show, tx)
	})

	// return the error if we got one
	if err != nil {
		return show, err
	}

	return show, nil
}

// Puts a show into the database, returning the previous show if it existed.
func (database *Database) PutShow(show Show) (Show, bool, error) {
	oldShow := Show{
		Day:  show.Day,
		Hour: show.Hour,
	}

	err := database.db.Update(func(tx *buntdb.Tx) error {
		// try and get the old show, ignoring a not found error
		err := FillShowFromTransaction(&oldShow, tx)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}

		showTime := show.Day + " " + show.Hour

		// set the new show
		_, _, err = tx.Set(ShowPrefix+showTime+HostSuffix, strings.Join(show.Hosts, ";"), nil)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(ShowPrefix+showTime+NameSuffix, show.Name, nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return oldShow, false, err
	}

	// now just return everything
	return oldShow, oldShow.Name != "", nil
}

// Deletes a show from the database.
func (database *Database) DeleteShow(day, hour string) error {
	err := database.db.Update(func(tx *buntdb.Tx) error {
		showTime := day + " " + hour

		_, err := tx.Delete(ShowPrefix + showTime + HostSuffix)
		if err != nil {
			return err
		}

		_, err = tx.Delete(ShowPrefix + showTime + NameSuffix)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// Obtains a show given a time, filling it in from a transaction.
func FillShowFromTransaction(show *Show, tx *buntdb.Tx) error {
	showTime := show.Day + " " + show.Hour

	logrus.WithField("show", showTime).Debug("Getting show from database...")

	host, err := tx.Get(ShowPrefix + showTime + HostSuffix)
	if err != nil {
		return err
	}

	name, err := tx.Get(ShowPrefix + showTime + NameSuffix)
	if err != nil {
		return err
	}

	// update the show
	show.Hosts = strings.Split(host, ";")
	show.Name = name

	return nil
}
