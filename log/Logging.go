// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

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

func printf(reqdLevel Level, prefix string, format string, v ...interface{}) {
	if level <= reqdLevel {
		fmt.Printf(prefix+format+"\n", v...)
	}
}

func sysprintf(reqdLevel Level, prefix string, format string, v ...interface{}) {
	if level <= reqdLevel {
		systemlog.Printf(prefix+format+"\n", v...)
	}
}

// Trace outputs a formatted tracing message
func Trace(format string, v ...interface{}) {
	printf(TRACE, "", format, v...)
}

// Debug outputs a formatted debugging message
func Debug(format string, v ...interface{}) {
	printf(DEBUG, "", format, v...)
}

// Info outputs a formatted information message
func Info(format string, v ...interface{}) {
	printf(INFO, "", format, v...)
}

// Warn outputs a formatted warning message
func Warn(format string, v ...interface{}) {
	printf(WARN, "Warning: ", format, v...)
}

// Error outputs a formatted error message
func Error(format string, v ...interface{}) {
	sysprintf(ERROR, "Error: ", format, v...)
}

// Fatal outputs a formatted fatal error message
func Fatal(format string, v ...interface{}) {
	// fatal errors go to the system fatal error handler
	if level <= FATAL {
		systemlog.Fatalf("Fatal: "+format+"\n", v...)
	}

	// if logLevel is NONE (the only level higher than FATAL) we terminate anyway but without a message
	os.Exit(1)
}
