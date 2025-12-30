package main

import (
	"fmt"
)

type Level int

const (
    LevelDebug Level = -4
    LevelInfo  Level = 0
    LevelWarn  Level = 4
    LevelError Level = 8
)

var LogLevel = LevelInfo

func logger(msg string, level Level) {
	if level >= LogLevel {
		fmt.Printf("%s %s\n", logLevelToString(level), msg)
	}
}

func logLevelToString(level Level) string {
	switch level {
	case LevelDebug:
		return "[DEBUG]"
	case LevelInfo:
		return "[INFO]"
	case LevelWarn:
		return "[WARN]"
	case LevelError:
		return "[ERROR]"
	default:
		return "[UNKNOWN]"
	}
}

func Debug(msg string) { logger(msg, LevelDebug) }
func Info(msg string)  { logger(msg, LevelInfo) }
func Warn(msg string)  { logger(msg, LevelWarn) }
func Err(msg string)   { logger(msg, LevelError) }
