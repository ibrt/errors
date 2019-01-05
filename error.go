// Package errors extends the functionality of Go's built-in error interface: it attaches stack traces to errors and
// supports behaviors such as carrying debug values and HTTP status codes. Additional behaviors that store metadata on
// errors can be easily implemented by users. Multiple errors can be merged into a compound one.
//
// The package provides several built-in behaviors (Prefix, Metadata, Callers, Skip, PublicMessage, HTTPStatus), ways to
// wrap and create errors (Errorf, MustErrorf, (Maybe)?Wrap, (Maybe)?MustWrap, (Maybe)?WrapRecover), ways to compound
// errors ((Maybe)?Append, ((Maybe?)Split) and utilities (Assert, Ignore, IgnoreClose, Unwrap, Equals).
//
// A wrapped error augments Go built-in errors with stack traces and additional behaviors. It can be created from an
// existing error using one of the Wrap function variants, or from scratch using one of the Errorf variants. To clients
// it appears to be a generic Go error, but functions in this library understand its magic and can manipulate it
// accordingly.
//
// This library also supports compound errors, i.e. an error composed by multiple inner errors. They can be created
// using one of the Append function variants, and - if needed - decomposed back using one of the Split function
// variants. Compound errors can generally be consumed as any other error, although they are subject to special
// treatment within this library as documented on individual methods.
package errors

import (
	"fmt"
	"io"
	"strings"
)

type wrappedError struct {
	err      error
	metadata map[interface{}]interface{}
}

// Error implements error.
func (e *wrappedError) Error() string {
	return GetPrefix(e) + e.err.Error()
}

type wrappedErrors []*wrappedError

// Error implements error.
func (e wrappedErrors) Error() string {
	b := strings.Builder{}
	b.WriteString("multiple errors: ")

	for i, err := range e {
		if i > 0 {
			b.WriteString(" · ")
		}
		b.WriteString(err.Error())
	}

	return b.String()
}

// Wrap wraps the given error, applying the given behaviors plus Callers. If the given error is already wrapped, only
// the behaviors are applied. If the given error is a compound error, Wrap is applied to the last inner error.
func Wrap(err error, behaviors ...Behavior) error {
	if err == nil {
		panic("nil error")
	}

	behaviors = append([]Behavior{Callers(), Skip(2)}, behaviors...)

	if wErr, ok := err.(*wrappedError); ok {
		Behaviors(behaviors...)(true, wErr)
		return wErr
	}

	if wErrs, ok := err.(wrappedErrors); ok {
		Behaviors(behaviors...)(true, wErrs[len(wErrs)-1])
		return wErrs
	}

	wErr := &wrappedError{
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
	case *wrappedError:
		return r
	case wrappedErrors:
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

// Append appends newErr to existingErr, creating or extending a compound error. All parameters can be unwrapped errors,
// wrapped errors, or compound errors. If newErr is a compound error, all the inner errors are appended.
//
// If existingErr is  nil, Append behaves like Wrap on newErr, thus returning a non-compound error. In all other cases a
// compound error is returned.
func Append(existingErr, newErr error) error {
	if newErr == nil {
		panic("nil error")
	}

	if existingErr == nil {
		return Wrap(newErr)
	}

	var wErrs wrappedErrors

	switch err := existingErr.(type) {
	case *wrappedError:
		wErrs = wrappedErrors{err}
	case wrappedErrors:
		wErrs = err
	default:
		wErrs = wrappedErrors{Wrap(err).(*wrappedError)}
	}

	switch err := newErr.(type) {
	case *wrappedError:
		return append(wErrs, err)
	case wrappedErrors:
		return append(wErrs, err...)
	default:
		return append(wErrs, Wrap(err).(*wrappedError))
	}
}

// MaybeAppend is like Append, but returns existingErr if newErr is nil.
func MaybeAppend(existingErr, newErr error) error {
	if newErr == nil {
		return existingErr
	}
	return Append(existingErr, newErr)
}

// Split allows access to the inner errors of a compound error. If the given error is not a compound error, the returned
// slice will contain such error as the single element.
func Split(err error) []error {
	if err == nil {
		panic("nil error")
	}

	switch err := err.(type) {
	case *wrappedError:
		return []error{err}
	case wrappedErrors:
		errs := make([]error, len(err))
		for i, err := range err {
			errs[i] = err
		}
		return errs
	default:
		return []error{err}
	}
}

// MaybeSplit is like Split, but returns an empty nil if err is nil.
func MaybeSplit(err error) []error {
	if err == nil {
		return nil
	}
	return Split(err)
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

// Unwrap undoes Wrap, returning the original error. If the given error is already unwrapped, it is simply returned
// as is. If the given error is a compound error, the last inner error is unwrapped and returned.
func Unwrap(err error) error {
	if wErr, ok := err.(*wrappedError); ok {
		return wErr.err
	}
	if wErrs, ok := err.(wrappedErrors); ok {
		return Unwrap(wErrs[len(wErrs)-1])
	}

	return err
}

// Equals returns true if the given error equals any of the given causes. If the given error is a compound error, Equals
// returns true if any of the inner errors equals any of the given causes. Both the given error and causes are
// unwrapped before checking for equality.
func Equals(err error, causes ...error) bool {
	if wErrs, ok := err.(wrappedErrors); ok {
		for _, wErr := range wErrs {
			if Equals(wErr, causes...) {
				return true
			}
		}

		return false
	}

	err = Unwrap(err)

	for _, cause := range causes {
		if wErrs, ok := cause.(wrappedErrors); ok {
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
