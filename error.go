// Package errors extends the functionality of Go's built-in error interface: it attaches stack traces to errors and
// supports behaviors such as carrying debug values and HTTP status codes. Additional behaviors can be easily
// implemented by users. The provided *Error type implements error and can be used interchangeably with code that
// expects a regular error return.
//
// The package provides several built-in behaviors (Prefix, Metadata, Callers, Skip, PublicMessage, HTTPStatus), ways to
// wrap and create errors (Errorf, MustErrorf, (Maybe)?Wrap, (Maybe)?MustWrap, (Maybe)?WrapRecover), and utilities
// (Assert, Ignore, IgnoreClose, Unwrap, Equals).
package errors

import (
	"fmt"
	"io"
)

// Error augments Go built-in errors with stack traces and additional behaviors.
type Error struct {
	err      error
	metadata map[interface{}]interface{}
}

// Error implements error.
func (e *Error) Error() string {
	return GetPrefix(e) + e.err.Error()
}

// Wrap wraps the given error into a *Error, applying the given behaviors plus Callers.
func Wrap(err error, behaviors ...Behavior) error {
	if err == nil {
		panic("nil error")
	}

	behaviors = append([]Behavior{Callers(), Skip(2)}, behaviors...)

	if wErr, ok := err.(*Error); ok {
		Behaviors(behaviors...)(true, wErr)
		return wErr
	}

	if wErrs, ok := err.(Errors); ok {
		Behaviors(behaviors...)(true, wErrs[len(wErrs)-1])
		return wErrs
	}

	wErr := &Error{
		err:      err,
		metadata: make(map[interface{}]interface{}),
	}

	Behaviors(behaviors...)(false, wErr)
	return wErr
}

// MaybeWrap is like Wrap, but returns nil if called with a nil error.
func MaybeWrap(err error, behaviors ...Behavior) error {
	if err == nil {
		return nil
	}

	behaviors = append(behaviors, Skip(1))
	return Wrap(err, behaviors...)
}

// MustWrap is like Wrap, but panics if the given error is non-nil.
func MustWrap(err error, behaviors ...Behavior) {
	if err == nil {
		panic("nil error")
	}

	behaviors = append(behaviors, Skip(1))
	panic(Wrap(err, behaviors...))
}

// MaybeMustWrap is like MustWrap, but does nothing if called with a nil error.
func MaybeMustWrap(err error, behaviors ...Behavior) {
	if err == nil {
		return
	}

	behaviors = append(behaviors, Skip(1))
	MustWrap(err, behaviors...)
}

// WrapRecover takes a recovered interface{} and converts it to a wrapped error.
func WrapRecover(r interface{}, behaviors ...Behavior) error {
	if r == nil {
		panic("nil recover")
	}

	behaviors = append(behaviors, Skip(1))

	switch r := r.(type) {
	case *Error:
		return r
	case Errors:
		return r
	case error:
		return Wrap(r, behaviors...)
	default:
		return Wrap(fmt.Errorf("%v", r), behaviors...)
	}
}

// MaybeWrapRecover is like WrapRecover but returns nil if called with a nil recover.
func MaybeWrapRecover(r interface{}, behaviors ...Behavior) error {
	if r == nil {
		return nil
	}

	behaviors = append(behaviors, Skip(1))
	return WrapRecover(r, behaviors...)
}

// Errorf formats a new error and wraps it.
// Note: arguments implementing Behavior are applied on wrapping, the others are passed to fmt.Errorf().
func Errorf(format string, behaviorOrArg ...interface{}) error {
	behaviors := make([]Behavior, 0, len(behaviorOrArg))
	args := make([]interface{}, 0, len(behaviorOrArg))

	for _, behaviorOrArg := range behaviorOrArg {
		if behavior, ok := behaviorOrArg.(Behavior); ok {
			behaviors = append(behaviors, behavior)
		} else {
			args = append(args, behaviorOrArg)
		}
	}

	behaviors = append(behaviors, Skip(1))
	return Wrap(fmt.Errorf(format, args...), behaviors...)
}

// MustErrorf is like Errorf but panics instead of returning the error.
func MustErrorf(format string, behaviorOrArg ...interface{}) {
	behaviorOrArg = append(behaviorOrArg, Skip(1))
	panic(Errorf(format, behaviorOrArg...))
}

// Append appends newErr to existingErr, returning a compound error. See the Errors type GoDoc for more information
// about compound errors. All parameters can be of type error, *Error, or Errors. If existingErr is nil, Append()
// behaves like Wrap() on newErr. If newErr is of type Errors, all the contained errors are appended.
func Append(existingErr, newErr error) error {
	if newErr == nil {
		panic("nil error")
	}

	if existingErr == nil {
		return Wrap(newErr)
	}

	var wErrs Errors

	switch err := existingErr.(type) {
	case *Error:
		wErrs = Errors{err}
	case Errors:
		wErrs = err
	default:
		wErrs = Errors{Wrap(err).(*Error)}
	}

	switch err := newErr.(type) {
	case *Error:
		return append(wErrs, err)
	case Errors:
		return append(wErrs, err...)
	default:
		return append(wErrs, Wrap(err).(*Error))
	}
}

// Assert is like MustErrorf if cond is false, does nothing otherwise.
func Assert(cond bool, format string, behaviorOrArg ...interface{}) {
	if cond {
		return
	}

	behaviorOrArg = append(behaviorOrArg, Skip(1))
	MustErrorf(format, behaviorOrArg...)
}

// Ignore does nothing. It is used in cases where the error is intentionally ignored to suppress lint errors.
func Ignore(_ error) {
	// intentionally empty
}

// IgnoreClose calls Close on the given io.Closer, ignoring the returned error. Handy for the defer Close pattern.
func IgnoreClose(c io.Closer) {
	Ignore(c.Close())
}

// Unwrap returns the wrapped error if the given error is of type *Error, the last wrapped error if the given error is
// of type Errors, or the given error itself otherwise.
func Unwrap(err error) error {
	if wErr, ok := err.(*Error); ok {
		return wErr.err
	}
	if wErrs, ok := err.(Errors); ok {
		return Unwrap(wErrs[len(wErrs)-1])
	}

	return err
}

// Equals returns true if the given error equals any of the given "causes". If the given error is a compound error, it
// returns true if any of the individual errors equals any of the given "causes".
func Equals(err error, causes ...error) bool {
	if wErrs, ok := err.(Errors); ok {
		for _, wErr := range wErrs {
			if Equals(wErr, causes...) {
				return true
			}
		}

		return false
	}

	err = Unwrap(err)

	for _, cause := range causes {
		if wErrs, ok := cause.(Errors); ok {
			for _, cause := range wErrs {
				if err == Unwrap(cause) {
					return true
				}
			}
		} else {
			if err == Unwrap(cause) {
				return true
			}
		}
	}

	return false
}
