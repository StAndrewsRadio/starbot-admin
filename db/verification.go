package db

import (
	"fmt"
	"time"

	"github.com/tidwall/buntdb"
)

var (
	verificationOptions = &buntdb.SetOptions{
		Expires: true,
		TTL:     15 * time.Minute,
	}
)

// Checks if a verification code for a user is valid.
func (database *Database) CheckVerification(userID, code string) (bool, error) {
	tx, err := database.db.Begin(false)
	if err != nil {
		return false, err
	}

	idToVerify, err := tx.Get(fmt.Sprintf(VerificationCodeToUser, code))
	if err == buntdb.ErrNotFound {
		return false, tx.Rollback()
	} else if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	return idToVerify == userID, tx.Rollback()
}

// Stores a verification code for 15 minutes.
func (database *Database) StoreVerificationCode(userID, email, code string) error {
	tx, err := database.db.Begin(true)
	if err != nil {
		return err
	}

	_, _, err = tx.Set(fmt.Sprintf(VerificationCodeToUser, code), userID, verificationOptions)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, _, err = tx.Set(fmt.Sprintf(VerificationCodeToEmail, code), email, verificationOptions)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Checks if an email is registered to a user, returning the user ID if it is.
func (database *Database) IsEmailRegistered(email string) (bool, string, error) {
	tx, err := database.db.Begin(false)
	if err != nil {
		return false, "", err
	}

	userID, err := tx.Get(fmt.Sprintf(EmailToUser, email))
	_ = tx.Rollback()
	if err == buntdb.ErrNotFound {
		return false, "", nil
	} else if err != nil {
		return false, "", err
	} else {
		return true, userID, nil
	}
}

// Validates a user, storing their email in the emails database for future validation.
func (database *Database) ValidateUser(userID, code string) error {
	tx, err := database.db.Begin(true)
	if err != nil {
		return err
	}

	email, err := tx.Get(fmt.Sprintf(VerificationCodeToEmail, code))
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, _, err = tx.Set(fmt.Sprintf(EmailToUser, email), userID, nil)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Removes an email address from the database, returning the user ID of the user who held validation if it was found.
func (database *Database) InvalidateEmail(email string) (string, error) {
	tx, err := database.db.Begin(true)
	if err != nil {
		return "", err
	}

	userID, err := tx.Delete(fmt.Sprintf(EmailToUser, email))
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}

	return userID, tx.Commit()
}
