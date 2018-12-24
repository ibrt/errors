package errors_test

import (
	"fmt"
	"net/http"
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

func ExampleBehaviors() {
	doSomething := func() error {
		behaviors := errors.Behaviors(errors.Prefix("prefix"), errors.HTTPStatus(http.StatusInternalServerError))
		return errors.Errorf("test error", behaviors)
	}

	if err := doSomething(); err != nil {
		fmt.Println(err.Error())
		fmt.Println(errors.GetHTTPStatus(err))
	}

	// Output:
	// prefix: test error
	// 500
}

func TestBehaviors(t *testing.T) {
	facets := errors.Behaviors(errors.Prefix("other error"), errors.Prefix("next error"))
	err := errors.Errorf("test error", facets)
	require.Equal(t, "next error: other error: test error", err.Error())
}
