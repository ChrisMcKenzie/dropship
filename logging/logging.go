package logging

import "github.com/Sirupsen/logrus"

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.DebugLevel
}

func GetLogger() *logrus.Logger {
	return log
}
