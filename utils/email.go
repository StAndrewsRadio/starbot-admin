package utils

import (
	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"gopkg.in/gomail.v2"
)

type Emailer struct {
	dialer  *gomail.Dialer
	address string
}

// Creates a new emailer.
func NewEmailer(config *cfg.Config) *Emailer {
	return &Emailer{
		dialer: gomail.NewDialer(config.GetString(cfg.EmailDomain), config.GetInt(cfg.EmailPort),
			config.GetString(cfg.EmailAddress), config.GetString(cfg.EmailPassword)),
		address: config.GetString(cfg.EmailAddress),
	}
}
