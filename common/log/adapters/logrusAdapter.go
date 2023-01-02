package adapters

import (
	"github.com/sirupsen/logrus"
)

var (
	logLevels  = []string{"debug", "info", "error"}
	logFormats = []string{"json", "text"}
)

type utcFormatter struct {
	f logrus.Formatter
}

func (f *utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return f.f.Format(e)
}

func NewLogger() *logrus.Logger {
	l := logrus.New()

	l.SetLevel(logrus.InfoLevel)
	l.SetFormatter(&utcFormatter{
		f: &logrus.JSONFormatter{},
	})

	return l
}
