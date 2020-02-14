//See https://github.com/aszxqw/simplelog/blob/master/README.md
package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

const LOG_FORMAT = "%5s %s %12s:%d\t\t%s\n"

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
	"\033[1;36mINFO\033[0m",
	"\033[1;33mWARN\033[0m",
	"\033[1;31mERROR\033[0m",
	"FATAL",
}

var startTime time.Time = time.Now()
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
	timeDiff := time.Since(startTime)
	_, filename, line, _ := runtime.Caller(2)
	_, filename = path.Split(filename)
	// TODO print timeDiff with fixed precision.
	fmt.Printf(LOG_FORMAT, LEVEL_STRINGS[level], timeDiff, filename, line, message)
}
