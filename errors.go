package errors

import "strings"

// Errors describes a compound error, i.e. an error composed by multiple errors. It can be created using Append(). It
// can generally be consumed as any other error (or *Error), although it is subject to some special casings:
//
// Error() returns "multiple errors: " followed by the concatenation of all individual Error() calls.  On Wrap(),
// Behaviors are always applied to the last error. Metadata getters such as GetPublicMessage() search starting from the
// last error and return the first match found. Equals() matches against all of the errors, returns true if at least one
// match is found. Unwrap() unwraps and returns the last error. An empty Errors is treated as nil (e.g. in MaybeWrap()).
type Errors []*Error

// Error implements error. It returns "multiple errors: " followed by the concatenation of all individual Error() calls.
func (e Errors) Error() string {
	if e == nil || len(e) == 0 {
		panic("nil error")
	}

	b := strings.Builder{}
	b.WriteString("multiple errors: ")

	for i, err := range e {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(err.Error())
	}

	return b.String()
}
