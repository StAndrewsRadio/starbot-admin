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
