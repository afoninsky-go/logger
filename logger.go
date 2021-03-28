package logger

import (
	"net/http"
	"time"

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

// Middleware returns http logging middleware
func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		l.WithContext(r.Context()).
			WithField("method", r.Method).
			WithField("duration", time.Since(start)).
			Info(r.URL)
	})
}
