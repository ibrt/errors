package errors_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/ibrt/errors"
	"github.com/stretchr/testify/require"
)

func TestFormatCallers(t *testing.T) {
	callers := make([]uintptr, 1024)
	formattedCallers := errors.FormatCallers(callers[:runtime.Callers(1, callers[:])])

	require.Len(t, formattedCallers, 3)
	require.True(t, strings.HasPrefix(formattedCallers[0], "errors_test.TestFormatCallers"))
	require.True(t, strings.HasPrefix(formattedCallers[1], "testing.tRunner"))
	require.True(t, strings.HasPrefix(formattedCallers[2], "runtime.goexit"))
}
