//See https://github.com/aszxqw/simplelog/blob/master/README.md
package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

const LOG_FORMAT = "%5s %8dms %12s:%-4d\t%s\n"

// Log Level: DEBUG INFO WARN ERROR FATAL
const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

var LEVEL_STRINGS_COLOR = [...]string{
	"DEBUG",
	"\033[1;36mINFO\033[0m",
	"\033[1;33mWARN\033[0m",
	"\033[1;31mERROR\033[0m",
	"FATAL",
}

var LEVEL_STRINGS = [...]string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

var startTime time.Time = time.Now()
var log_level = LEVEL_DEBUG
var color_enabled = true

// If SetLevel is not called , log_level is the default value : LEVEL_DEBUG (the lowest level)
func SetLogLevel(level int) {
	log_level = level
}

// DisableColors make sure no ASCII sequence is used, useful for browser console logs
func DisableColors() {
	color_enabled = false
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
	timeDiff := time.Since(startTime).Milliseconds()
	_, filename, line, _ := runtime.Caller(2)
	_, filename = path.Split(filename)
	if color_enabled {
		fmt.Printf(LOG_FORMAT, LEVEL_STRINGS_COLOR[level], timeDiff, filename, line, message)
	} else {
		fmt.Printf(LOG_FORMAT, LEVEL_STRINGS[level], timeDiff, filename, line, message)
	}
}
