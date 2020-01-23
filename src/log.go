//See https://github.com/aszxqw/simplelog/blob/master/README.md
package main

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

// Log Format : [log_level] [time] [filename]:[line] [message]
const FORMAT = "%s %s %s:%d\t\t%s\n"

// Log Level: DEBUG INFO WARN ERROR FATAL
const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

var LEVEL_STRINGS = [...]string{
	"DEBUG",
	" INFO",
	" WARN",
	"ERROR",
	"FATAL",
}

var log_level = LEVEL_DEBUG

// If SetLevel is not called , log_level is the default value : LEVEL_DEBUG (the lowest level)
func SetLogLevel(level int) {
	log_level = level
}

func Debug(format string, v ...interface{}) {
	logByLevel(LEVEL_DEBUG, fmt.Sprintf(format, v...))
}

func Info(format string, v ...interface{}) {
	logByLevel(LEVEL_INFO, fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	logByLevel(LEVEL_WARN, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	logByLevel(LEVEL_ERROR, fmt.Sprintf(format, v...))
}

func Fatal(format string, v ...interface{}) {
	logByLevel(LEVEL_FATAL, fmt.Sprintf(format, v...))
}

func logByLevel(level int, message string) {
	if level < log_level {
		return
	}
	time_str := fmt.Sprintf("%s", time.Now())[:19]
	_, filename, line, _ := runtime.Caller(2)
	_, filename = path.Split(filename)
	fmt.Printf(FORMAT, LEVEL_STRINGS[level], time_str, filename, line, message)
}
