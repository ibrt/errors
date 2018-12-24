package errors_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ibrt/errors"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	status        int
	constant      errors.Behavior
	publicMessage string
}

var testCases = []testCase{
	{http.StatusBadRequest, errors.HTTPStatusBadRequest, "bad-request"},
	{http.StatusUnauthorized, errors.HTTPStatusUnauthorized, "unauthorized"},
	{http.StatusPaymentRequired, errors.HTTPStatusPaymentRequired, "payment-required"},
	{http.StatusForbidden, errors.HTTPStatusForbidden, "forbidden"},
	{http.StatusNotFound, errors.HTTPStatusNotFound, "not-found"},
	{http.StatusMethodNotAllowed, errors.HTTPStatusMethodNotAllowed, "method-not-allowed"},
	{http.StatusNotAcceptable, errors.HTTPStatusNotAcceptable, "not-acceptable"},
	{http.StatusProxyAuthRequired, errors.HTTPStatusProxyAuthRequired, "proxy-auth-required"},
	{http.StatusRequestTimeout, errors.HTTPStatusRequestTimeout, "request-timeout"},
	{http.StatusConflict, errors.HTTPStatusConflict, "conflict"},
	{http.StatusGone, errors.HTTPStatusGone, "gone"},
	{http.StatusLengthRequired, errors.HTTPStatusLengthRequired, "length-required"},
	{http.StatusPreconditionFailed, errors.HTTPStatusPreconditionFailed, "precondition-failed"},
	{http.StatusRequestEntityTooLarge, errors.HTTPStatusRequestEntityTooLarge, "request-entity-too-large"},
	{http.StatusRequestURITooLong, errors.HTTPStatusRequestURITooLong, "request-uri-too-long"},
	{http.StatusUnsupportedMediaType, errors.HTTPStatusUnsupportedMediaType, "unsupported-media-type"},
	{http.StatusRequestedRangeNotSatisfiable, errors.HTTPStatusRequestedRangeNotSatisfiable, "requested-range-not-satisfiable"},
	{http.StatusExpectationFailed, errors.HTTPStatusExpectationFailed, "expectation-failed"},
	{http.StatusTeapot, errors.HTTPStatusTeapot, "i-am-a-teapot"},
	{http.StatusMisdirectedRequest, errors.HTTPStatusMisdirectedRequest, "misdirected-request"},
	{http.StatusUnprocessableEntity, errors.HTTPStatusUnprocessableEntity, "unprocessable-entity"},
	{http.StatusLocked, errors.HTTPStatusLocked, "locked"},
	{http.StatusFailedDependency, errors.HTTPStatusFailedDependency, "failed-dependency"},
	{http.StatusUpgradeRequired, errors.HTTPStatusUpgradeRequired, "upgrade-required"},
	{http.StatusPreconditionRequired, errors.HTTPStatusPreconditionRequired, "precondition-required"},
	{http.StatusTooManyRequests, errors.HTTPStatusTooManyRequests, "too-many-requests"},
	{http.StatusRequestHeaderFieldsTooLarge, errors.HTTPStatusRequestHeaderFieldsTooLarge, "request-header-fields-too-large"},
	{http.StatusUnavailableForLegalReasons, errors.HTTPStatusUnavailableForLegalReasons, "unavailable-for-legal-reasons"},
	{http.StatusInternalServerError, errors.HTTPStatusInternalServerError, "internal-server-error"},
	{http.StatusNotImplemented, errors.HTTPStatusNotImplemented, "not-implemented"},
	{http.StatusBadGateway, errors.HTTPStatusBadGateway, "bad-gateway"},
	{http.StatusServiceUnavailable, errors.HTTPStatusServiceUnavailable, "service-unavailable"},
	{http.StatusGatewayTimeout, errors.HTTPStatusGatewayTimeout, "gateway-timeout"},
	{http.StatusHTTPVersionNotSupported, errors.HTTPStatusHTTPVersionNotSupported, "http-version-not-supported"},
	{http.StatusVariantAlsoNegotiates, errors.HTTPStatusVariantAlsoNegotiates, "variant-also-negotiates"},
	{http.StatusInsufficientStorage, errors.HTTPStatusInsufficientStorage, "insufficient-storage"},
	{http.StatusLoopDetected, errors.HTTPStatusLoopDetected, "loop-detected"},
	{http.StatusNotExtended, errors.HTTPStatusNotExtended, "not-extended"},
	{http.StatusNetworkAuthenticationRequired, errors.HTTPStatusNetworkAuthenticationRequired, "network-authentication-required"},
}

func TestHTTPStatusConstants(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(http.StatusText(testCase.status), func(t *testing.T) {
			require.Equal(t, testCase.status, errors.GetHTTPStatus(errors.Errorf("err", testCase.constant)))
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

func ExampleHTTPPublicMessage() {
	doSomething := func() error {
		return errors.Errorf("test error", errors.HTTPPublicMessage(http.StatusInternalServerError))
	}

	if err := doSomething(); err != nil {
		fmt.Println(errors.GetPublicMessage(err))
	}

	// Output:
	// internal-server-error
}

func ExampleHTTPPublicMessage_unknown() {
	doSomething := func() error {
		return errors.Errorf("test error", errors.HTTPPublicMessage(499))
	}

	if err := doSomething(); err != nil {
		fmt.Println(errors.GetPublicMessage(err))
	}

	// Output:
	// unknown
}

func TestHTTPPublicMessage(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(http.StatusText(testCase.status), func(t *testing.T) {
			err := errors.Errorf("err", errors.HTTPPublicMessage(testCase.status))
			require.Equal(t, testCase.publicMessage, errors.GetPublicMessage(err))
		})
	}

	t.Run("Unknown", func(t *testing.T) {
		err := errors.Errorf("err", errors.HTTPPublicMessage(499))
		require.Equal(t, "unknown", errors.GetPublicMessage(err))
	})
}

func ExampleHTTPError() {
	doSomething := func() error {
		return errors.Errorf("test error", errors.HTTPError(http.StatusInternalServerError))
	}

	if err := doSomething(); err != nil {
		fmt.Println(errors.GetHTTPStatus(err))
		fmt.Println(errors.GetPublicMessage(err))
	}

	// Output:
	// 500
	// internal-server-error
}

func TestHTTPError(t *testing.T) {
	err := errors.Errorf("test error", errors.HTTPError(http.StatusInternalServerError))
	require.Equal(t, http.StatusInternalServerError, errors.GetHTTPStatus(err))
	require.Equal(t, "internal-server-error", errors.GetPublicMessage(err))
}
