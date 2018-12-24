# errors [![Build Status](https://travis-ci.org/ibrt/errors.svg?branch=master)](https://travis-ci.org/ibrt/errors) [![Go Report Card](https://goreportcard.com/badge/github.com/ibrt/errors)](https://goreportcard.com/report/github.com/ibrt/errors) [![Test Coverage](https://codecov.io/gh/ibrt/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/ibrt/errors) [![Go Docs](https://godoc.org/github.com/ibrt/errors?status.svg)](http://godoc.org/github.com/ibrt/errors)

Package `errors` extends the functionality of Go's built-in error interface: it attaches stack traces to errors and supports behaviors such as carrying debug values and HTTP status codes. Additional behaviors can be easily implemented by users. The provided `*Error` type implements `error` and can be used interchangeably with code that expects a regular error return.

The package provides several built-in behaviors (`Prefix`, `Metadata`, `Callers`, `Skip`, `PublicMessage`, `HTTPStatus`, `HTTPPublicMessage`, `HTTPError`), ways to wrap and create errors (`Errorf`, `MustErrorf`, `(Maybe)?Wrap`, `(Maybe)?MustWrap`, `(Maybe)?WrapRecover`), and utilities (`Assert`, `Ignore`, `IgnoreClose`, `Unwrap`, `Equals`). See [GoDoc](https://godoc.org/github.com/ibrt/errors) for detailed usage examples.

#### Basic Example

```go
func doSomething() error {
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
    fmt.Println(errors.FormatCallers(errors.GetCallers(err)))
}
```

```
read failed: EOF
500
internal server error
true
true
[ some_pkg.SomeFunc (/Users/../some_pkg/file.go:26) ... ]
```
