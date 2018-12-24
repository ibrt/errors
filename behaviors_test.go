package errors_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ibrt/errors"
	"github.com/stretchr/testify/require"
)

func ExampleMetadata() {
	doSomething := func() error {
		return errors.Errorf("test error", errors.Metadata("key", "value"))
	}

	if err := doSomething(); err != nil {
		fmt.Println(errors.GetMetadata(err, "key"))
	}

	// Output:
	// value
}

func MyValue(value string) errors.Behavior {
	return errors.Metadata(reflect.TypeOf(MyValue), value)
}

func GetMyValue(err error) string {
	if value, ok := errors.GetMetadata(err, reflect.TypeOf(MyValue)).(string); ok {
		return value
	}
	return ""
}

func ExampleMetadata_customBehavior() {
	//  func MyValue(value string) errors.Behavior {
	//    return errors.Metadata(reflect.TypeOf(MyValue), value)
	//  }
	//
	//  func GetMyValue(err error) string {
	//    if value, ok := errors.GetMetadata(err, reflect.TypeOf(MyValue)).(string); ok {
	//      return value
	//    }
	//    return ""
	//  }

	doSomething := func() error {
		return errors.Errorf("test error", MyValue("my value"))
	}

	if err := doSomething(); err != nil {
		fmt.Println(GetMyValue(err))
	}

	// Output:
	// my value
}

func TestMetadata(t *testing.T) {
	require.Nil(t, errors.GetMetadata(fmt.Errorf("test error"), "key"))
	err := errors.Errorf("test error", errors.Metadata("key", "value"))
	require.Equal(t, "value", errors.GetMetadata(err, "key"))
}

func TestSkip(t *testing.T) {
	err := errors.Errorf("test error")
	require.True(t, strings.HasPrefix(errors.FormatCallers(errors.GetCallers(err))[0], "errors_test.TestSkip"))
	err = errors.Errorf("test error", errors.Skip(1))
	require.True(t, strings.HasPrefix(errors.FormatCallers(errors.GetCallers(err))[0], "testing.tRunner"))
	err = errors.Wrap(errors.Errorf("test error"), errors.Skip(1))
	require.True(t, strings.HasPrefix(errors.FormatCallers(errors.GetCallers(err))[0], "errors_test.TestSkip"))
}

func ExamplePrefix() {
	doSomething := func() error {
		return errors.Errorf("test error", errors.Prefix("prefix"))
	}

	if err := doSomething(); err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// prefix: test error
}

func TestPrefix(t *testing.T) {
	err := errors.Errorf("test error", errors.Prefix("other error"), errors.Prefix("next error"))
	require.Equal(t, "next error: other error: test error", err.Error())
	err = errors.Wrap(err, errors.Prefix("final error"))
	require.Equal(t, "final error: next error: other error: test error", err.Error())
}

func ExamplePublicMessage() {
	doSomething := func() error {
		return errors.Errorf("a detailed error", errors.PublicMessage("a public error"))
	}

	if err := doSomething(); err != nil {
		fmt.Println(err.Error())
		fmt.Println(errors.GetPublicMessage(err))
	}

	// Output:
	// a detailed error
	// a public error
}

func ExamplePublicMessage_default() {
	doSomething := func() error {
		return errors.Errorf("a detailed error")
	}

	if err := doSomething(); err != nil {
		fmt.Println(errors.GetPublicMessageOrDefault(err, "default"))
	}

	// Output:
	// default
}

func TestPublicMessage(t *testing.T) {
	err := errors.Errorf("test error")
	require.Equal(t, "", errors.GetPublicMessage(err))
	require.Equal(t, "default", errors.GetPublicMessageOrDefault(err, "default"))
	err = errors.Errorf("test error", errors.PublicMessage("public message"))
	require.Equal(t, "public message", errors.GetPublicMessage(err))
	require.Equal(t, "public message", errors.GetPublicMessageOrDefault(err, "default"))
	err = errors.Wrap(err, errors.PublicMessage("another public message"))
	require.Equal(t, "another public message", errors.GetPublicMessage(err))
}
