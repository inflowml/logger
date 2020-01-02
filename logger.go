// Package logger offers utility functions related to logging.
// The logger.go module implements generic logging functions.
package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	// Terminal colour codes.
	reset   = "\033[0m"
	red     = "\033[31m"
	orange  = "\033[33m"
	blue    = "\033[36m"
	grey    = "\033[37m"
	green   = "\033[32m"
	magenta = "\033[95m"
)

// Enabled dictates whether any log messages are displayed.
var Enabled = true

// getTime returns a time.Time object that represents the current time.
var getTime = time.Now

// halt causes the current execution thread to stop.
var halt = func(v interface{}) { panic(v) }

// Debug displays a debug message to stdout only if env-var DEBUG is 1
// The given message and arguments are used in a similar fashion to printf().
func Debug(msg string, args ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		print("DEBUG", msg, green, args...)
	}
}

// Network displays a network message to stdout.
// The given message and arguments are used in a similar fashion to printf().
func Network(msg string, args ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		print("NET", msg, grey, args...)
	}
}

// SQL displays a SQL message to stdout.
// The given message and arguments are used in a similar fashion to printf().
func SQL(msg string, args ...interface{}) {
	if os.Getenv("SQL_LOG") != "" {
		print("SQL", msg, magenta, args...)
	}
}

// Info displays an information message to stdout.
// The given message and arguments are used in a similar fashion to printf().
func Info(msg string, args ...interface{}) {
	print("INFO", msg, blue, args...)
}

// Warning displays a warning message to stdout.
// The given message and arguments are used in a similar fashion to printf().
func Warning(msg string, args ...interface{}) {
	print("WARN", msg, orange, args...)
}

// Error displays an error message to stdout.
// The given message and arguments are used in a similar fashion to printf().
func Error(msg string, args ...interface{}) {
	print("ERROR", msg, red, args...)
}

// Fatal displays an error message to stdout and calls the exit() function.
// The given message and arguments are used in a similar fashion to printf().
func Fatal(msg string, args ...interface{}) {
	Error(msg, args...)
	halt("logger.Fatal() was called.")
}

// print prepends a timestamp and the given severity to the provided message.
// The given message and arguments are used in a similar fashion to printf().
func print(sev, msg, clr string, args ...interface{}) {
	if Enabled {
		now := getTime().Format("2006/01/02 15:04:05.000")
		prefix := fmt.Sprintf("[%s] %s%s%s: ", now, clr, sev, reset)
		format := prefix + msg + "\n"
		fmt.Printf(format, args...)
	}
}
