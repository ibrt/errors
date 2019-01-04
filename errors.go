package errors

import "strings"

// Errors describes a compound error. It can be created using Append(). Behaviors on Wrap are always applied to the
// last error in the compound error, while getters cycle through the errors starting from the last one, and return the
// first match found.
type Errors []*Error

// Error implements error. It returns a concatenation
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
