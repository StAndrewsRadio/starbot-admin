package program

import "github.com/sirupsen/logrus"

var (
	Commands []interface{}
)

// Registers all commands to the go-flags parser.
func RegisterCommands() {
	logrus.Debug("Registering commands...")

	for _, command := range Commands {
		command.(func())()
	}
}
