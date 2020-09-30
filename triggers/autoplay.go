package triggers

import (
	"fmt"
	"net/http"

	"github.com/StAndrewsRadio/starbot-admin/jobs"
	"github.com/sirupsen/logrus"
)

func autoplay(writer http.ResponseWriter, request *http.Request) {
	logger := logrus.WithField("trigger", "autoplay")

	_, err := fmt.Fprint(writer, "Autoplay job scheduled!")

	go jobs.StartAutoplay(botSession, userSession, config, database, false, false)

	if err != nil {
		logger.WithError(err).Error("An error occurred.")
	} else {
		logger.Debug("Autoplay requested.")
	}
}
