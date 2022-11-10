package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

type LogField = logrus.Fields

func InitLogger(useJSONFormatter bool) {
	if useJSONFormatter {
		// Log as JSON instead of the default ASCII formatter.
		Log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{})
	}
	// Output to stdout instead of the default stderr
	Log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	Log.SetLevel(logrus.DebugLevel)
}
