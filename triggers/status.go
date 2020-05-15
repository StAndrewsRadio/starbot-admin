package triggers

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func status(writer http.ResponseWriter, request *http.Request) {
	_, err := fmt.Fprint(writer, "all good here :~)")

	logger := logrus.WithField("trigger", "status")
	if err != nil {
		logger.WithError(err).Error("An error occurred whilst returning a status request.")
	} else {
		logger.Debug("Status update requested.")
	}
}
