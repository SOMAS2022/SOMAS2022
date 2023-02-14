package logging

import (
	"fmt"
	"infra/game/state"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var gameState *state.State

type LogField = logrus.Fields

type Level uint32

const (
	Trace = iota
	Debug
	Info
	Warn
	Error
)

func InitLogger(verbose, useJSONFormatter, debug bool, id string, state *state.State) {
	if useJSONFormatter {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
	}

	if verbose {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(io.Discard)
	}
	runID = id
	gameState = state

	if debug {
		log.SetLevel(logrus.TraceLevel)
	} else {
		LogToFile(Warn, nil, "'Trace' and 'Debug' messages hidden. Run with '-d' or 'make runDebug' to see these logs.", LevelStages{})
		log.SetLevel(logrus.InfoLevel)
	}
}

func Log(lvl Level, fields LogField, msg string) {
	switch lvl {
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

func CombineFields(a LogField, b LogField) LogField {
	for k, v := range b {
		if a[k] != nil {
			Log(Warn, nil, fmt.Sprintf("Overwriting '%s' default logging field to '%s'", k, v))
		}
		a[k] = v
	}
	return a
}
