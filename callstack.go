// Package logger offers utility functions related to logging.
// The timer.go module implements timing-related logging functions.
package logger

import (
	"errors"
	"fmt"
	"runtime"
)

// GetCallerName returns the name of the function that called the calling function.
func GetCallerName() (string, error) {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "", errors.New("failed to retrieve caller function PC")
	}

	caller := runtime.FuncForPC(pc)
	if caller == nil {
		return "", errors.New("failed to derive function name from PC")
	}
	return caller.Name(), nil
}

// GetCallStack returns a slice of strings representing the call stack starting
// from the calling function.
func GetCallStack(size int) ([]string, error) {
	stack := make([]string, 0, size)
	for i := 0; i < size; i++ {
		pc, _, _, ok := runtime.Caller(size - i)
		if !ok {
			return stack, errors.New("failed to retrieve caller function PC")
		}

		caller := runtime.FuncForPC(pc)
		if caller == nil {
			return stack, errors.New("failed to derive function name from PC")
		}

		file, line := caller.FileLine(pc)
		name := fmt.Sprintf("%s:%d", file, line)
		stack = append(stack, name)
	}
	return stack, nil
}
