package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Emailer struct {
	dialer           *gomail.Dialer
	address          string
	validEmails      []string
	verificationFile string
}

// Creates a new emailer.
func NewEmailer(config *cfg.Config) (*Emailer, error) {
	emailer := &Emailer{
		dialer: gomail.NewDialer(config.GetString(cfg.EmailDomain), int(config.GetInt(cfg.EmailPort)),
			config.GetString(cfg.EmailAddress), config.GetString(cfg.EmailPassword)),
		address:     config.GetString(cfg.EmailAddress),
		validEmails: []string{},
	}

	file, err := os.Open(config.GetString(cfg.VerificationAllowedEmails))
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		emailer.validEmails = append(emailer.validEmails, scanner.Text())
	}

	verificationEmailContents, err := ioutil.ReadFile(config.GetString(cfg.VerificationEmailContents))
	if err != nil {
		return nil, err
	}

	emailer.verificationFile = string(verificationEmailContents)

	return emailer, nil
}

// Checks if an email is within the list of valid emails.
func (emailer *Emailer) IsValidEmail(email string) bool {
	return StringSliceContains(emailer.validEmails, email)
}

// Sends an email to an address
func (emailer *Emailer) SendVerificationEmail(to, subject, prefix, code string) {
	message := gomail.NewMessage()
	message.SetHeader("From", emailer.address)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", fmt.Sprintf(emailer.verificationFile, prefix, code))

	err := emailer.dialer.DialAndSend(message)
	if err != nil {
		logrus.WithError(err).Error("An error occurred whilst sending an email.")
	}
}
