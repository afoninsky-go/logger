package logger

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrus.Entry
}

type MiddlewareConfig struct {
	IgnorePaths []string
}

func New() *Logger {
	var log = logrus.NewEntry(logrus.New())
	return &Logger{*log}
}

func NewSTDLogger() *Logger {
	log := New()
	log.Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	return log
}

func (l *Logger) FatalIfError(err error) {
	if err != nil {
		l.Fatal(err)
	}
}

// CreateMiddleware returns http logging middleware
func (l *Logger) CreateMiddleware(cfg *MiddlewareConfig) func(next http.Handler) http.Handler {
	ignoreRoutes := map[string]bool{}
	if cfg != nil {
		for _, p := range cfg.IgnorePaths {
			ignoreRoutes[p] = true
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			o := &responseObserver{ResponseWriter: w}
			next.ServeHTTP(o, r)

			if _, ok := ignoreRoutes[r.URL.Path]; ok {
				return
			}
			l.WithContext(r.Context()).
				WithField("method", r.Method).
				WithField("duration", time.Since(start).Milliseconds()).
				WithField("status", o.status).
				Info(r.URL)
		})
	}
}

// responseObserver is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}
