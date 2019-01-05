# errors [![Build Status](https://travis-ci.org/ibrt/errors.svg?branch=master)](https://travis-ci.org/ibrt/errors) [![Go Report Card](https://goreportcard.com/badge/github.com/ibrt/errors)](https://goreportcard.com/report/github.com/ibrt/errors) [![Test Coverage](https://codecov.io/gh/ibrt/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/ibrt/errors) [![Go Docs](https://godoc.org/github.com/ibrt/errors?status.svg)](http://godoc.org/github.com/ibrt/errors)

Package errors extends the functionality of Go's built-in error interface: it attaches stack traces to errors and
supports behaviors such as carrying debug values and HTTP status codes. Additional behaviors that store metadata on
errors can be easily implemented by users. Multiple errors can be merged into a compound one.

The package provides several built-in behaviors (`Prefix`, `Metadata`, `Callers`, `Skip`, `PublicMessage`, 
`HTTPStatus`), ways to wrap and create errors `((Must?)Errorf`, `(Maybe)?(Must)?Wrap`, `(Maybe)?(Must?)WrapRecover)`, 
ways to compound errors `((Maybe)?Append`, `((Maybe?)Split)` and utilities (`Assert`, `Ignore`, `IgnoreClose`, `Unwrap`,
`Equals`).

A wrapped error augments Go built-in errors with stack traces and additional behaviors. It can be created from an
existing error using one of the Wrap function variants, or from scratch using one of the Errorf variants. To clients it 
appears to be a generic Go error, but functions in this library understand its magic and can manipulate it accordingly.

This library also supports compound errors, i.e. an error composed by multiple inner errors. They can be created using 
one of the Append function variants, and - if needed - decomposed back using one of the Split function variants. 
Compound errors can generally be consumed as any other error, although they are subject to special treatment within this
library as documented on individual methods.

#### Example (Basic)

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
	fmt.Println(errors.getCallers(err))
}
```

Outputs:

```
read failed: EOF
500
internal server error
true
true
[ some_pkg.SomeFunc (/Users/../some_pkg/file.go:26) ... ]
```

#### Example (Compound)

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

func doSomethingElse() error {
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
```

Outputs

```
Output:
multiple errors: read failed: EOF · some error
500
internal server error
true
false
read failed: EOF
some error
```

See the [GoDoc](https://godoc.org/github.com/ibrt/errors) for detailed usage examples.
