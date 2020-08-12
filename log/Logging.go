package log

import (
	"fmt"
	systemlog "log"
	"os"
)

// Level is the level of the logger
type Level int

const (
	// TRACE represents simple tracing messages
	TRACE Level = 0

	// DEBUG represents messages used for debugging purposes
	DEBUG Level = 1

	// INFO represents general information messages
	INFO Level = 2

	// WARN represents warnings
	WARN Level = 3

	// ERROR represents errors
	ERROR Level = 4

	// FATAL represents fatal errors
	FATAL Level = 5

	// NONE is used to turn off all logging
	NONE Level = 6
)

func init() {
	systemlog.SetFlags(systemlog.Ldate)
	systemlog.SetFlags(systemlog.Ltime)
	systemlog.SetFlags(systemlog.Lshortfile)
}

var level = ERROR

// GetLevel retrieves the current logging level
func GetLevel() Level {
	return level
}

// SetLevel sets the current logging level
func SetLevel(newLevel Level) {
	level = newLevel
}

func print(reqdLevel Level, v ...interface{}) {
	if level <= reqdLevel {
		fmt.Println(v...)
	}
}

func printf(reqdLevel Level, format string, v ...interface{}) {
	if level <= reqdLevel {
		fmt.Printf(format, v...)
	}
}

func sysprint(reqdLevel Level, v ...interface{}) {
	if level <= reqdLevel {
		systemlog.Println(v...)
	}
}

func sysprintf(reqdLevel Level, format string, v ...interface{}) {
	if level <= reqdLevel {
		systemlog.Printf(format, v...)
	}
}

// Trace outputs a tracing message
func Trace(v ...interface{}) {
	print(TRACE, v...)
}

// Tracef outputs a formatted tracing message
func Tracef(format string, v ...interface{}) {
	printf(TRACE, format, v...)
}

// Debug outputs a debugging message
func Debug(v ...interface{}) {
	print(DEBUG, v...)
}

// Debugf outputs a formatted debugging message
func Debugf(format string, v ...interface{}) {
	printf(DEBUG, format, v...)
}

// Info outputs an information message
func Info(v ...interface{}) {
	print(INFO, v...)
}

// Infof outputs a formatted information message
func Infof(format string, v ...interface{}) {
	printf(INFO, format, v...)
}

// Warn outputs a warning message
func Warn(v ...interface{}) {
	sysprint(WARN, v...)
}

// Warnf outputs a formatted warning message
func Warnf(format string, v ...interface{}) {
	sysprintf(WARN, format, v...)
}

// Error outputs an error message
func Error(v ...interface{}) {
	sysprint(ERROR, v...)
}

// Errorf outputs a formatted error message
func Errorf(format string, v ...interface{}) {
	sysprintf(ERROR, format, v...)
}

// Fatal outputs a fatal error message
func Fatal(v ...interface{}) {
	// fatal errors go to the system fatal error handler
	if level <= FATAL {
		systemlog.Fatal(v...)
	}

	// if logLevel is NONE (the only level higher than FATAL) we terminate anyway but without a message
	os.Exit(1)
}

// Fatalf outputs a formatted fatal error message
func Fatalf(format string, v ...interface{}) {
	// fatal errors go to the system fatal error handler
	if level <= FATAL {
		systemlog.Fatalf(format, v...)
	}

	// if logLevel is NONE (the only level higher than FATAL) we terminate anyway but without a message
	os.Exit(1)
}
