// Package log provides a logger interface for logger libraries
// so that bee does not depend on any of them directly.
// It also includes a default implementation using Logrus (used by bee previously).
package log

import (
	"github.com/luminosita/bee/common/log/adapters"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

// Logger serves as an adapter interface for logger libraries
// so that bee does not depend on any of them directly.
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

var (
	once   sync.Once
	logger Logger
)

func GetLogger() Logger {
	once.Do(func() { // <-- atomic, does not allow repeating
		logger = adapters.NewLogger()
	})

	return logger
}

type utcFormatter struct {
	f logrus.Formatter
}

func (f *utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return f.f.Format(e)
}

func SetLogger(level string, format string) {
	var logLevel logrus.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = logrus.DebugLevel
	case "", "info":
		logLevel = logrus.InfoLevel
	case "error":
		logLevel = logrus.ErrorLevel
	default:
		logLevel = logrus.InfoLevel
	}

	var formatter utcFormatter
	switch strings.ToLower(format) {
	case "", "text":
		formatter.f = &logrus.TextFormatter{DisableColors: true}
	case "json":
		formatter.f = &logrus.JSONFormatter{}
	default:
		formatter.f = &logrus.JSONFormatter{}
	}

	GetLogger().(*logrus.Logger).SetLevel(logLevel)
	GetLogger().(*logrus.Logger).SetFormatter(&formatter)
}
