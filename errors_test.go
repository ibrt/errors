package errors_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ibrt/errors"
)

func ExampleErrors() {
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

	// Output:
	// multiple errors: read failed: EOF · some error
	// 500
	// internal server error
	// true
	// false
}
