package util

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type GeometricTraceLogger struct {
	numberOfEntries int64
}

func (l *GeometricTraceLogger) Msg(s string) {
	l.numberOfEntries++

	doLog := false
	if l.numberOfEntries <= 10 ||
		(l.numberOfEntries <= 100 && l.numberOfEntries%10 == 0) ||
		(l.numberOfEntries <= 1000 && l.numberOfEntries%100 == 0) ||
		(l.numberOfEntries <= 10000 && l.numberOfEntries%1000 == 0) ||
		(l.numberOfEntries <= 100000 && l.numberOfEntries%1000 == 0) ||
		(l.numberOfEntries <= 1000000 && l.numberOfEntries%100000 == 0) {
		doLog = true
	}

	if doLog {
		log.Trace().Msg(s)
	}
}

func (l *GeometricTraceLogger) MsgEvent(e *zerolog.Event, msg string) {
	l.numberOfEntries++

	doLog := false
	if l.numberOfEntries <= 10 ||
		(l.numberOfEntries <= 100 && l.numberOfEntries%10 == 0) ||
		(l.numberOfEntries <= 1000 && l.numberOfEntries%100 == 0) ||
		(l.numberOfEntries <= 10000 && l.numberOfEntries%1000 == 0) ||
		(l.numberOfEntries <= 100000 && l.numberOfEntries%1000 == 0) ||
		(l.numberOfEntries <= 1000000 && l.numberOfEntries%100000 == 0) {
		doLog = true
	}

	if doLog {
		if msg != "" {
			e.Msg(msg)
		} else {
			e.Send()
		}
	}
}
