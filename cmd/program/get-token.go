package program

import (
	"github.com/StAndrewsRadio/starbot-admin/vars"
	"github.com/sirupsen/logrus"
)

type GetTokenArgs struct {
	Email    string `positional-arg-name:"email" description:"The email of the account"`
	Password string `positional-arg-name:"password" description:"The password of the account"`
}

type GetTokenCommand struct {
	Arguments GetTokenArgs `positional-args:"true" required:"true"`
}

var getTokenCommand GetTokenCommand

func init() {
	Commands = append(Commands, func() {
		_, err := vars.Parser.AddCommand("get-token", "Gets a login token for a user account.",
			"Gets a token using the Discord API that is used to login a user account.",
			&getTokenCommand)
		if err != nil {
			logrus.WithError(err).Fatal("An error occurred whilst parsing the get-token command!")
		}
	})
}

func (cmd *GetTokenCommand) Execute(args []string) error {
	// init
	starbot, err := vars.InitialiseStarbot(
		vars.WithUserSession(cmd.Arguments.Email, cmd.Arguments.Password),
	)

	// check if any error happened
	if err != nil {
		return err
	}

	logrus.WithField("token", starbot.UserSession.Token).Info("Token obtained.")
	return nil
}
