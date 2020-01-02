package logger

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// RPC converts the given message into an error that suitable for return from an
// RPC function:
//
//   1. In a production environment, a generic error message is returned to avoid
//      leaking implementation details to public callers.
//
//   2. In a debug environment, the provided message is formatted and returned to
//      facilitate cross-microservice debugging.
//
// Either way, the message is logged.
func RPC(msg string, args ...interface{}) error {
	if len(msg) > 0 {
		display := strings.ToUpper(msg[:1]) + msg[1:] + "."
		Error(display, args...)
	}

	if os.Getenv("DEBUG") != "" {
		return fmt.Errorf(msg, args...)
	}
	return errors.New("failed to complete RPC request")
}
