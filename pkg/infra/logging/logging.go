package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

type LogField = logrus.Fields

type Level uint32

const (
	Trace = iota
	Debug
	Info
	Warn
	Error
)

func InitLogger(useJSONFormatter bool, debug bool) {
	if useJSONFormatter {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
	}

	log.SetOutput(os.Stdout)

	if debug {
		log.SetLevel(logrus.TraceLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
}

func Log(err Level, fields LogField, msg string) {
	switch err {
	case Trace:
		log.WithFields(fields).Trace(msg)
	case Debug:
		log.WithFields(fields).Debug(msg)
	case Info:
		log.WithFields(fields).Info(msg)
	case Warn:
		log.WithFields(fields).Warn(msg)
	case Error:
		log.WithFields(fields).Error(msg)
	default:
		log.WithFields(fields).Info(msg)
	}
}
