package triggers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/StAndrewsRadio/starbot-admin/cfg"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	accessPassword          string
	server                  *http.Server
	botSession, userSession *discordgo.Session
	config                  *cfg.Config
)

type Handler struct {
}

// Sets the triggers http server up
func SetupTriggers(bs, us *discordgo.Session, c *cfg.Config) {
	// store stuff
	botSession = bs
	userSession = us
	config = c
	accessPassword = config.GetString(cfg.TriggersPassword)

	server = &http.Server{Addr: config.GetString(cfg.TriggersAddress), Handler: Handler{}}
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			// we only care about the error if it's not the closing
			logrus.WithError(err).Error("An error occurred with the triggers HTTP server.")
		}
	}
}

// Closes the http server
func Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("An error occurred whilst closing the HTTP server.")
	}
}

// Handler for http server
func (Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logrus.WithField("url", request.URL.String()).Debug("Got triggers request...")

	// password check
	if !checkPassword(request) {
		_, _ = fmt.Fprint(writer, "Invalid password.")
	} else {
		// now let's redirect the request
		if val, ok := request.URL.Query()["trigger"]; ok {
			if len(val) == 1 {
				switch val[0] {
				case "autoplay":
					autoplay(writer, request)
				default:
					status(writer, request)
				}
			}
		} else {
			status(writer, request)
		}
	}
}

// Checks if a request contains the correct password
func checkPassword(request *http.Request) bool {
	if val, ok := request.URL.Query()["password"]; ok {
		if len(val) == 1 {
			return val[0] == accessPassword
		}
	}

	return false
}
