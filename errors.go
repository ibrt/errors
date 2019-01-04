package errors

import "strings"

// Errors describes a compound error, i.e. an error composed by multiple inner errors. It can be created using Append.
// It can generally be consumed as any other error (or *Error), although it is subject to some special casings:
//
// Error returns "multiple errors: " followed by the concatenation of all inner Error calls. On Wrap, Behaviors are
// always applied to the last inner error. Metadata getters such as GetPublicMessage search starting from the last inner
// error and return the first match found. Equals matches against all of the inner errors, returns true if at least one
// matches. Unwrap unwraps and returns the last inner error. An empty Errors is generally treated as a nil error.
type Errors []*Error

// Error implements error. It returns "multiple errors: " followed by the concatenation of all individual Error calls.
func (e Errors) Error() string {
	if e == nil || len(e) == 0 {
		panic("nil error")
	}

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
