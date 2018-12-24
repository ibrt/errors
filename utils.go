package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// FormatCallers returns a human-readable version of callers.
func FormatCallers(callers []uintptr) []string {
	frames := runtime.CallersFrames(callers)
	formattedCallers := make([]string, 0, len(callers))

	for {
		frame, more := frames.Next()

		formattedCallers = append(
			formattedCallers,
			fmt.Sprintf("%v (%v:%v)", filepath.Base(frame.Function), frame.File, frame.Line))

		if !more {
			break
		}
	}

	return formattedCallers
}

// Behaviors compounds multiple behaviors in a single Behavior.
func Behaviors(behaviors ...Behavior) Behavior {
	return func(doubleWrap bool, e *Error) {
		for _, behavior := range behaviors {
			behavior(doubleWrap, e)
		}
	}
}
