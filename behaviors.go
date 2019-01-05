package errors

import (
	"fmt"
	"reflect"
	"runtime"
)

// Behavior describes an additional behavior to be applied to the error.
type Behavior func(doubleWrap bool, err error)

// Metadata returns a Behavior that stores the given key/value pair in the error metadata.
func Metadata(key, value interface{}) Behavior {
	return func(_ bool, err error) {
		err.(*wrappedError).metadata[key] = value
	}
}

// GetMetadata extracts the given key from the error metadata.
// If the given error is compound, the key is searched starting from the last inner error, and the first match (if any)
// is returned.
func GetMetadata(err error, key interface{}) interface{} {
	if e, ok := err.(*wrappedError); ok {
		return e.metadata[key]
	}

	if e, ok := err.(wrappedErrors); ok {
		for i := len(e) - 1; i >= 0; i-- {
			if v, ok := e[i].metadata[key]; ok {
				return v
			}
		}
	}

	return nil
}

// Callers is a Behavior that stores a stack trace in the error metadata.
// It is automatically applied on Wrap.
func Callers() Behavior {
	return func(doubleWrap bool, err error) {
		if GetCallers(err) == nil {
			callers := make([]uintptr, 1024)
			Metadata(reflect.ValueOf(Callers), callers[:runtime.Callers(2, callers[:])])(doubleWrap, err)
		}
	}
}

// GetCallers extracts a stack trace from the error metadata, if any.
// It returns nil if no stack trace was set. The callers behavior is automatically applied on wrap.
func GetCallers(err error) []uintptr {
	if callers, ok := GetMetadata(err, reflect.ValueOf(Callers)).([]uintptr); ok {
		return callers
	}
	return nil
}

// GetCallersOrCurrent extracts a stack trace from the error metadata, if any.
// It returns the current stack trace if no stack trace was set.
func GetCallersOrCurrent(err error) []uintptr {
	if callers := GetCallers(err); callers != nil {
		return callers
	}
	callers := make([]uintptr, 1024)
	return callers[:runtime.Callers(2, callers[:])]
}

// Skip returns a Behavior that skips the given amount of trailing frames in the stack trace.
func Skip(skip int) Behavior {
	return func(doubleWrap bool, err error) {
		if callers := GetCallers(err); !doubleWrap && callers != nil && len(callers) > skip {
			Metadata(reflect.ValueOf(Callers), callers[skip:])(doubleWrap, err)
		}
	}
}

// Prefix returns a Behavior that prepends a prefix to the error message.
// The prefixFormat and parameters are first passed through fmt.Sprintf().
func Prefix(prefixFormat string, a ...interface{}) Behavior {
	return func(doubleWrap bool, err error) {
		Metadata(reflect.ValueOf(Prefix), fmt.Sprintf(prefixFormat, a...)+": "+GetPrefix(err))(doubleWrap, err)
	}
}

// GetPrefix returns the computed error prefix on the error, if any.
// It returns "" if no prefix was set.
func GetPrefix(err error) string {
	if prefix, ok := GetMetadata(err, reflect.ValueOf(Prefix)).(string); ok {
		return prefix
	}
	return ""
}

// PublicMessage returns a behavior that stores a public message in the error metadata.
// It is useful in API servers where detailed errors are logged, while a different message is returned to clients.
func PublicMessage(message string) Behavior {
	return Metadata(reflect.ValueOf(PublicMessage), message)
}

// GetPublicMessage extracts a public message from the error metadata, if any.
// It returns "" if no public message was set.
func GetPublicMessage(err error) string {
	if message, ok := GetMetadata(err, reflect.ValueOf(PublicMessage)).(string); ok {
		return message
	}
	return ""
}

// GetPublicMessageOrDefault extracts a public message from the error metadata, if any.
// It returns the given default public message if no public message was set.
func GetPublicMessageOrDefault(err error, defaultMessage string) string {
	if message := GetPublicMessage(err); message != "" {
		return message
	}
	return defaultMessage
}
