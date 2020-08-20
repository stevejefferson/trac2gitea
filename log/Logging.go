// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	systemlog "log"
	"os"
	"runtime/debug"
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
	// systemlog.SetFlags(systemlog.Ldate)
	// systemlog.SetFlags(systemlog.Ltime)
	// systemlog.SetFlags(systemlog.Lshortfile)
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

func prependArg(arg interface{}, existing []interface{}) []interface{} {
	var newArgs []interface{}
	newArgs = append(newArgs, arg)
	newArgs = append(newArgs, existing)
	return newArgs
}

func print(reqdLevel Level, prefix string, v ...interface{}) {
	if level <= reqdLevel {
		fmt.Println(prependArg(prefix, v)...)
	}
}

func printf(reqdLevel Level, prefix string, format string, v ...interface{}) {
	if level <= reqdLevel {
		fmt.Printf(prefix+format, v...)
	}
}

func sysprint(reqdLevel Level, prefix string, v ...interface{}) {
	if level <= reqdLevel {
		systemlog.Println(prependArg(prefix, v)...)
		debug.PrintStack()
	}
}

func sysprintf(reqdLevel Level, prefix string, format string, v ...interface{}) {
	if level <= reqdLevel {
		systemlog.Printf(prefix+format, v...)
		debug.PrintStack()
	}
}

// Trace outputs a tracing message
func Trace(v ...interface{}) {
	print(TRACE, "", v...)
}

// Tracef outputs a formatted tracing message
func Tracef(format string, v ...interface{}) {
	printf(TRACE, "", format, v...)
}

// Debug outputs a debugging message
func Debug(v ...interface{}) {
	print(DEBUG, "", v...)
}

// Debugf outputs a formatted debugging message
func Debugf(format string, v ...interface{}) {
	printf(DEBUG, "", format, v...)
}

// Info outputs an information message
func Info(v ...interface{}) {
	print(INFO, "", v...)
}

// Infof outputs a formatted information message
func Infof(format string, v ...interface{}) {
	printf(INFO, "", format, v...)
}

// Warn outputs a warning message
func Warn(v ...interface{}) {
	print(WARN, "Warning: ", v...)
}

// Warnf outputs a formatted warning message
func Warnf(format string, v ...interface{}) {
	printf(WARN, "Warning: ", format, v...)
}

// Error outputs an error message
func Error(v ...interface{}) {
	sysprint(ERROR, "Error: ", v...)
}

// Errorf outputs a formatted error message
func Errorf(format string, v ...interface{}) {
	sysprintf(ERROR, "Error: ", format, v...)
}

// Fatal outputs a fatal error message
func Fatal(v ...interface{}) {
	// fatal errors go to the system fatal error handler
	if level <= FATAL {
		debug.PrintStack()
		systemlog.Fatal(prependArg("Fatal: ", v)...)
	}

	// if logLevel is NONE (the only level higher than FATAL) we terminate anyway but without a message
	os.Exit(1)
}

// Fatalf outputs a formatted fatal error message
func Fatalf(format string, v ...interface{}) {
	// fatal errors go to the system fatal error handler
	if level <= FATAL {
		debug.PrintStack()
		systemlog.Fatalf("Fatal: "+format, v...)
	}

	// if logLevel is NONE (the only level higher than FATAL) we terminate anyway but without a message
	os.Exit(1)
}
