package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = &logrus.Logger{
	Out:       os.Stdout,
	Formatter: &logrus.TextFormatter{DisableTimestamp: true},
	Level:     logrus.InfoLevel,
}
