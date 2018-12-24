package errors_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ibrt/errors"
	"github.com/stretchr/testify/require"
)

var testCases = map[int]errors.Behavior{
	http.StatusBadRequest:                    errors.HTTPStatusBadRequest,
	http.StatusUnauthorized:                  errors.HTTPStatusUnauthorized,
	http.StatusPaymentRequired:               errors.HTTPStatusPaymentRequired,
	http.StatusForbidden:                     errors.HTTPStatusForbidden,
	http.StatusNotFound:                      errors.HTTPStatusNotFound,
	http.StatusMethodNotAllowed:              errors.HTTPStatusMethodNotAllowed,
	http.StatusNotAcceptable:                 errors.HTTPStatusNotAcceptable,
	http.StatusProxyAuthRequired:             errors.HTTPStatusProxyAuthRequired,
	http.StatusRequestTimeout:                errors.HTTPStatusRequestTimeout,
	http.StatusConflict:                      errors.HTTPStatusConflict,
	http.StatusGone:                          errors.HTTPStatusGone,
	http.StatusLengthRequired:                errors.HTTPStatusLengthRequired,
	http.StatusPreconditionFailed:            errors.HTTPStatusPreconditionFailed,
	http.StatusRequestEntityTooLarge:         errors.HTTPStatusRequestEntityTooLarge,
	http.StatusRequestURITooLong:             errors.HTTPStatusRequestURITooLong,
	http.StatusUnsupportedMediaType:          errors.HTTPStatusUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable:  errors.HTTPStatusRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed:             errors.HTTPStatusExpectationFailed,
	http.StatusTeapot:                        errors.HTTPStatusTeapot,
	http.StatusUnprocessableEntity:           errors.HTTPStatusUnprocessableEntity,
	http.StatusLocked:                        errors.HTTPStatusLocked,
	http.StatusFailedDependency:              errors.HTTPStatusFailedDependency,
	http.StatusUpgradeRequired:               errors.HTTPStatusUpgradeRequired,
	http.StatusPreconditionRequired:          errors.HTTPStatusPreconditionRequired,
	http.StatusTooManyRequests:               errors.HTTPStatusTooManyRequests,
	http.StatusRequestHeaderFieldsTooLarge:   errors.HTTPStatusRequestHeaderFieldsTooLarge,
	http.StatusUnavailableForLegalReasons:    errors.HTTPStatusUnavailableForLegalReasons,
	http.StatusInternalServerError:           errors.HTTPStatusInternalServerError,
	http.StatusNotImplemented:                errors.HTTPStatusNotImplemented,
	http.StatusBadGateway:                    errors.HTTPStatusBadGateway,
	http.StatusServiceUnavailable:            errors.HTTPStatusServiceUnavailable,
	http.StatusGatewayTimeout:                errors.HTTPStatusGatewayTimeout,
	http.StatusHTTPVersionNotSupported:       errors.HTTPStatusHTTPVersionNotSupported,
	http.StatusVariantAlsoNegotiates:         errors.HTTPStatusVariantAlsoNegotiates,
	http.StatusInsufficientStorage:           errors.HTTPStatusInsufficientStorage,
	http.StatusLoopDetected:                  errors.HTTPStatusLoopDetected,
	http.StatusNotExtended:                   errors.HTTPStatusNotExtended,
	http.StatusNetworkAuthenticationRequired: errors.HTTPStatusNetworkAuthenticationRequired,
}

func TestHTTPStatusConstants(t *testing.T) {
	for status, constant := range testCases {
		t.Run(http.StatusText(status), func(t *testing.T) {
			require.Equal(t, status, errors.GetHTTPStatus(errors.Errorf("err", constant)))
		})
	}
}

func ExampleHTTPStatus() {
	doSomething := func() error {
		return errors.Errorf("test error", errors.HTTPStatus(http.StatusInternalServerError))
	}

	if err := doSomething(); err != nil {
		fmt.Println(errors.GetHTTPStatus(err))
	}

	// Output:
	// 500
}

func ExampleHTTPStatus_default() {
	doSomething := func() error {
		return errors.Errorf("test error")
	}

	if err := doSomething(); err != nil {
		fmt.Println(errors.GetHTTPStatusOrDefault(err, http.StatusInternalServerError))
	}

	// Output:
	// 500
}

func TestHTTPStatus(t *testing.T) {
	err := errors.Errorf("test error")
	require.Equal(t, 0, errors.GetHTTPStatus(err))
	require.Equal(t, 200, errors.GetHTTPStatusOrDefault(err, http.StatusOK))
	err = errors.Errorf("test error", errors.HTTPStatus(http.StatusOK))
	require.Equal(t, http.StatusOK, errors.GetHTTPStatus(err))
	require.Equal(t, http.StatusOK, errors.GetHTTPStatusOrDefault(err, http.StatusInternalServerError))
	err = errors.Wrap(err, errors.HTTPStatus(http.StatusInternalServerError))
	require.Equal(t, http.StatusInternalServerError, errors.GetHTTPStatus(err))
}
