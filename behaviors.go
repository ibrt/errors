package errors

import (
	"reflect"
	"runtime"
)

// Behavior describes an additional behavior to be applied to the Error.
type Behavior func(bool, *Error)

// Metadata returns a Behavior that stores the given key/value pair in the error metadata.
func Metadata(key, value interface{}) Behavior {
	return func(_ bool, e *Error) {
		e.metadata[key] = value
	}
}

// GetMetadata extracts the given key from the error metadata.
func GetMetadata(err error, key interface{}) interface{} {
	if e, ok := err.(*Error); ok {
		return e.metadata[key]
	}
	return nil
}

// Callers is a Behavior that stores a stack trace in the error metadata.
// It is automatically applied on Wrap.
func Callers() Behavior {
	return func(doubleWrap bool, e *Error) {
		if GetCallers(e) == nil {
			callers := make([]uintptr, 1024)
			Metadata(reflect.ValueOf(Callers), callers[:runtime.Callers(2, callers[:])])(doubleWrap, e)
		}
	}
}

// GetCallers extracts a stack trace from the error metadata, if any.
func GetCallers(err error) []uintptr {
	if callers, ok := GetMetadata(err, reflect.ValueOf(Callers)).([]uintptr); ok {
		return callers
	}
	return nil
}

// Skip returns a Behavior that skips the given amount of trailing frames in the stack trace.
func Skip(skip int) Behavior {
	return func(doubleWrap bool, e *Error) {
		if callers := GetCallers(e); !doubleWrap && callers != nil && len(callers) > skip {
			Metadata(reflect.ValueOf(Callers), callers[skip:])(doubleWrap, e)
		}
	}
}

// Prefix returns a Behavior that prepends a prefix to the error message.
func Prefix(prefix string) Behavior {
	return func(doubleWrap bool, e *Error) {
		Metadata(reflect.ValueOf(Prefix), prefix+": "+GetPrefix(e))(doubleWrap, e)
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

// HTTPStatus returns a behavior that stores a HTTP status in the error metadata.
func HTTPStatus(status int) Behavior {
	return Metadata(reflect.ValueOf(HTTPStatus), status)
}

// GetHTTPStatus extracts a HTTP status from the error metadata, if any.
// It returns 0 if no HTTP status was set.
func GetHTTPStatus(err error) int {
	if status, ok := GetMetadata(err, reflect.ValueOf(HTTPStatus)).(int); ok {
		return status
	}
	return 0
}

// GetHTTPStatusOrDefault extracts a HTTP status from the error metadata, if any.
// It returns the given default HTTP status if no HTTP status was set.
func GetHTTPStatusOrDefault(err error, defaultStatus int) int {
	if status := GetHTTPStatus(err); status != 0 {
		return status
	}
	return defaultStatus
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

// Behaviors compounds multiple behaviors in a single Behavior.
func Behaviors(behaviors ...Behavior) Behavior {
	return func(doubleWrap bool, e *Error) {
		for _, behavior := range behaviors {
			behavior(doubleWrap, e)
		}
	}
}
