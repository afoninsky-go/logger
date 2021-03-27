package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrus.Entry
}

func NewSTDLogger() *Logger {
	var log = logrus.NewEntry(logrus.New())

	log.Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	return &Logger{*log}
}

func (l *Logger) FatalIfError(err error) {
	if err != nil {
		l.Fatal(err)
	}
}
