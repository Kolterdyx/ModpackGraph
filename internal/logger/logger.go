package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// Init initializes the global logger
func Init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// GetLogger returns the global logger
func GetLogger() *logrus.Logger {
	if Log == nil {
		Init()
	}
	return Log
}
