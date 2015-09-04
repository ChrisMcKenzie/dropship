package logging

import "github.com/Sirupsen/logrus"

var log = logrus.New()

func GetLogger() *logrus.Logger {
	return log
}
