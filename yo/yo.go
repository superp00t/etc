package yo

import (
	"time"
)

type LogLevel uint8

const (
	LOkay  LogLevel = 0
	LDebug LogLevel = 1
	LWarn  LogLevel = 2
	LFatal LogLevel = 3
)

type Logger interface {
	Log(LogData)
}

var stdLogger = NewConsoleColor()

func Attach(l Logger) {
	stdLogger = l
}

type LogData struct {
	Level LogLevel
	Time  time.Time
	Data  string
}
