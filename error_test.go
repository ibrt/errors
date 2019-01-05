package errors_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ibrt/errors"
	"github.com/stretchr/testify/require"
)

func Example_simple() {
	doSomething := func() error {
		if _, err := strings.NewReader("").Read(make([]byte, 1024)); err != nil {
			return errors.Wrap(err,
				errors.Prefix("read failed"),
				errors.HTTPStatus(http.StatusInternalServerError),
				errors.PublicMessage("internal server error"))
		}
		return nil
	}

	if err := doSomething(); err != nil {
		fmt.Println(err.Error())
		fmt.Println(errors.GetHTTPStatus(err))
		fmt.Println(errors.GetPublicMessage(err))
		fmt.Println(errors.Equals(err, io.EOF))
		fmt.Println(errors.Unwrap(err) == io.EOF)
	}

	// Output:
	// read failed: EOF
	// 500
	// internal server error
	// true
	// true
}

func Example_compound() {
	doSomething := func() error {
		if _, err := strings.NewReader("").Read(make([]byte, 1024)); err != nil {
			return errors.Wrap(err,
				errors.Prefix("read failed"),
				errors.HTTPStatus(http.StatusInternalServerError),
				errors.PublicMessage("internal server error"))
		}
		return nil
	}

	doSomethingElse := func() error {
		return fmt.Errorf("some error")
	}

	var errs error

	if err := doSomething(); err != nil {
		errs = errors.Append(errs, err)
	}
	if err := doSomethingElse(); err != nil {
		errs = errors.Append(errs, err)
	}

	if errs != nil {
		fmt.Println(errs.Error())
		fmt.Println(errors.GetHTTPStatus(errs))
		fmt.Println(errors.GetPublicMessage(errs))
		fmt.Println(errors.Equals(errs, io.EOF))
		fmt.Println(errors.Unwrap(errs) == io.EOF) // Unwrap returns the last error
	}

	for _, err := range errors.Split(errs) {
		fmt.Println(err.Error())
	}

	// Output:
	// multiple errors: read failed: EOF · some error
	// 500
	// internal server error
	// true
	// false
	// read failed: EOF
	// some error
}

func ExampleWrap() {
	doSomething := func() error {
		if _, err := strings.NewReader("").Read(make([]byte, 1024)); err != nil {
			return errors.Wrap(err,
				errors.Prefix("read failed"),
				errors.HTTPStatus(http.StatusInternalServerError))
		}
		return nil
	}

	if err := doSomething(); err != nil {
		fmt.Println(err.Error())
		fmt.Println(errors.GetHTTPStatus(err))
		fmt.Println(errors.Equals(err, io.EOF))
		fmt.Println(errors.Unwrap(err) == io.EOF)
	}

	// Output:
	// read failed: EOF
	// 500
	// true
	// true
}

func TestWrap(t *testing.T) {
	err := errors.Wrap(fmt.Errorf("test error"))
	require.Equal(t, "test error", err.Error())
	require.True(t, strings.HasPrefix(errors.FormatCallers(errors.GetCallers(err))[0], "errors_test.TestWrap"))
	require.PanicsWithValue(t, "nil error", func() { errors.Wrap(nil) })
}

func ExampleMaybeWrap() {
	doSomething := func() error {
		_, err := strings.NewReader("string").Read(make([]byte, 6))
		return errors.MaybeWrap(err)
	}

	if err := doSomething(); err == nil {
		fmt.Println("success")
	}

	// Output:
	// success
}

func TestMaybeWrap(t *testing.T) {
	err := errors.MaybeWrap(fmt.Errorf("test error"))
	require.Equal(t, "test error", err.Error())
	require.True(t, strings.HasPrefix(errors.FormatCallers(errors.GetCallers(err))[0], "errors_test.TestMaybeWrap"))
	require.Nil(t, errors.MaybeWrap(nil))
}

func ExampleMustWrap() {
	defer func() {
		fmt.Println("panic:", recover().(error).Error())
	}()

	doSomething := func() error {
		if _, err := strings.NewReader("").Read(make([]byte, 1024)); err != nil {
			errors.MustWrap(err, errors.Prefix("read failed"))
		}
		return nil
	}

	doSomething()

	// Output:
	// panic: read failed: EOF
}

func TestMustWrap(t *testing.T) {
	require.Panics(t, func() { errors.MustWrap(fmt.Errorf("test error")) })
	require.PanicsWithValue(t, "nil error", func() { errors.MustWrap(nil) })
}

func ExampleMaybeMustWrap() {
	doSomething := func() error {
		_, err := strings.NewReader("string").Read(make([]byte, 6))
		errors.MaybeMustWrap(err)
		return nil
	}

	if err := doSomething(); err == nil {
		fmt.Println("success")
	}

	// Output:
	// success
}

func TestMaybeMustWrap(t *testing.T) {
	require.Panics(t, func() { errors.MaybeMustWrap(fmt.Errorf("test error")) })
	require.NotPanics(t, func() { errors.MaybeMustWrap(nil) })
}

func ExampleWrapRecover() {
	defer func() {
		fmt.Println("panic:", errors.WrapRecover(recover(), errors.Prefix("read failed")).Error())
	}()

	panic("test error")

	// Output:
	// panic: read failed: test error
}

func TestWrapRecover(t *testing.T) {
	err := errors.WrapRecover("test error")
	require.Equal(t, "test error", err.Error())
	err = errors.WrapRecover(fmt.Errorf("test error"))
	require.Equal(t, "test error", err.Error())
	err = errors.WrapRecover(errors.Errorf("test error"))
	require.Equal(t, "test error", err.Error())
	require.PanicsWithValue(t, "nil recover", func() { errors.WrapRecover(nil) })
}

func ExampleMaybeWrapRecover() {
	defer func() {
		fmt.Println(errors.MaybeWrapRecover(recover()))
	}()

	fmt.Println("success")

	// Output:
	// success
	// <nil>
}

func TestMaybeWrapRecover(t *testing.T) {
	err := errors.MaybeWrapRecover("test error")
	require.Equal(t, "test error", err.Error())
	err = errors.MaybeWrapRecover(fmt.Errorf("test error"))
	require.Equal(t, "test error", err.Error())
	err = errors.MaybeWrapRecover(errors.Errorf("test error"))
	require.Equal(t, "test error", err.Error())
	require.Nil(t, errors.MaybeWrapRecover(nil))
}

func ExampleErrorf() {
	doSomething := func() error {
		return errors.Errorf("test error: %v", "EOF", errors.Prefix("prefix"))
	}

	if err := doSomething(); err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// prefix: test error: EOF
}

func TestErrorf(t *testing.T) {
	err := errors.Errorf("test error")
	require.Equal(t, "test error", err.Error())
	require.True(t, strings.HasPrefix(errors.FormatCallers(errors.GetCallers(err))[0], "errors_test.TestErrorf"))
	err = errors.Errorf("format %s", errors.Prefix("prefix"), "xxx", errors.HTTPStatus(http.StatusOK))
	require.Equal(t, "prefix: format xxx", err.Error())
	require.Equal(t, http.StatusOK, errors.GetHTTPStatus(err))
}

func ExampleMustErrorf() {
	defer func() {
		fmt.Println("panic:", errors.WrapRecover(recover()).Error())
	}()

	errors.MustErrorf("test error: %v", "EOF", errors.Prefix("prefix"))

	// Output:
	// panic: prefix: test error: EOF
}

func TestMustErrorf(t *testing.T) {
	require.Panics(t, func() { errors.MustErrorf("test error") })
}

func ExampleAssert() {
	defer func() {
		fmt.Println("panic:", errors.WrapRecover(recover()).Error())
	}()

	errors.Assert(true, "test error: %v", "true", errors.Prefix("prefix"))
	errors.Assert(false, "test error: %v", "false", errors.Prefix("prefix"))

	// Output:
	// panic: prefix: test error: false
}

func TestAssert(t *testing.T) {
	require.NotPanics(t, func() { errors.Assert(true, "test error") })
	require.Panics(t, func() { errors.Assert(false, "test error") })
}

func ExampleIgnore() {
	errors.Ignore(fmt.Errorf("test error"))

	// Output:
}

func TestIgnore(t *testing.T) {
	require.NotPanics(t, func() { errors.Ignore(fmt.Errorf("test")) })
}

type testCloser struct {
	closed bool
}

// Close implements io.Closer.
func (c *testCloser) Close() error {
	c.closed = true
	return nil
}

func ExampleIgnoreClose() {
	// type testCloser struct {
	//   closed bool
	// }
	//
	// func (c *testCloser) Close() error {
	//	 c.closed = true
	//	 return nil
	// }

	tc := &testCloser{}
	fmt.Println(tc.closed)
	errors.IgnoreClose(tc)
	fmt.Println(tc.closed)

	// Output:
	// false
	// true
}

func TestIgnoreClose(t *testing.T) {
	tc := &testCloser{}
	require.False(t, tc.closed)
	errors.IgnoreClose(tc)
	require.True(t, tc.closed)
}

func ExampleUnwrap() {
	fmt.Println(errors.Unwrap(nil))
	err := fmt.Errorf("test error")
	ret := errors.Unwrap(err)
	fmt.Println(ret == err)
	ret = errors.Unwrap(errors.Wrap(err))
	fmt.Println(ret == err)

	// Output:
	// <nil>
	// true
	// true
}

func TestUnwrap(t *testing.T) {
	require.Nil(t, errors.Unwrap(nil))
	err := fmt.Errorf("test error")
	ret := errors.Unwrap(err)
	require.Equal(t, ret, err)
	ret = errors.Unwrap(errors.Wrap(err))
	require.Equal(t, ret, err)
}

func ExampleEquals() {
	err := fmt.Errorf("test error")

	fmt.Println(errors.Equals(err, err))
	fmt.Println(errors.Equals(err, fmt.Errorf("other error"), err))
	fmt.Println(errors.Equals(err, fmt.Errorf("other error")))
	fmt.Println(errors.Equals(errors.Wrap(err), err))
	fmt.Println(errors.Equals(err, errors.Wrap(err)))
	fmt.Println(errors.Equals(errors.Wrap(err), errors.Wrap(err)))

	// Output:
	// true
	// true
	// false
	// true
	// true
	// true
}

func TestEquals(t *testing.T) {
	err := fmt.Errorf("test error")
	require.True(t, errors.Equals(err, err))
	require.True(t, errors.Equals(err, fmt.Errorf("other error"), err))
	require.False(t, errors.Equals(err, fmt.Errorf("other error")))
	require.True(t, errors.Equals(errors.Wrap(err), err))
	require.True(t, errors.Equals(err, errors.Wrap(err)))
	require.True(t, errors.Equals(errors.Wrap(err), errors.Wrap(err)))
	wErr := errors.Wrap(err)
	require.True(t, errors.Equals(wErr, wErr))
}
