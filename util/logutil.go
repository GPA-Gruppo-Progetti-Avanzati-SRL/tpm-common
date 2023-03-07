package util

import (
	"github.com/rs/zerolog"
)

type GeometricTraceLogger struct {
	numberOfEntries int64
}

func _isEnabled(numLines int64) bool {

	doLog := false
	if numLines <= 10 ||
		(numLines <= 100 && numLines%10 == 0) ||
		(numLines <= 1000 && numLines%100 == 0) ||
		(numLines <= 10000 && numLines%1000 == 0) ||
		(numLines <= 100000 && numLines%1000 == 0) ||
		(numLines <= 1000000 && numLines%100000 == 0) {
		doLog = true
	}

	return doLog
}

func (l *GeometricTraceLogger) IsEnabled() bool {
	return _isEnabled(l.numberOfEntries + 1)
}

func (l *GeometricTraceLogger) LogEvent(e *zerolog.Event, msg string) {
	l.numberOfEntries++

	if _isEnabled(l.numberOfEntries) {
		if msg != "" {
			e.Msg(msg)
		} else {
			e.Send()
		}
	}
}
