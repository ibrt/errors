package errors

import (
	"net/http"
	"reflect"
)

// Shorthand HTTPStatus behaviors for 4xx and 5xx HTTP status codes registered with IANA.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
var (
	HTTPStatusBadRequest                    = HTTPStatus(http.StatusBadRequest)
	HTTPStatusUnauthorized                  = HTTPStatus(http.StatusUnauthorized)
	HTTPStatusPaymentRequired               = HTTPStatus(http.StatusPaymentRequired)
	HTTPStatusForbidden                     = HTTPStatus(http.StatusForbidden)
	HTTPStatusNotFound                      = HTTPStatus(http.StatusNotFound)
	HTTPStatusMethodNotAllowed              = HTTPStatus(http.StatusMethodNotAllowed)
	HTTPStatusNotAcceptable                 = HTTPStatus(http.StatusNotAcceptable)
	HTTPStatusProxyAuthRequired             = HTTPStatus(http.StatusProxyAuthRequired)
	HTTPStatusRequestTimeout                = HTTPStatus(http.StatusRequestTimeout)
	HTTPStatusConflict                      = HTTPStatus(http.StatusConflict)
	HTTPStatusGone                          = HTTPStatus(http.StatusGone)
	HTTPStatusLengthRequired                = HTTPStatus(http.StatusLengthRequired)
	HTTPStatusPreconditionFailed            = HTTPStatus(http.StatusPreconditionFailed)
	HTTPStatusRequestEntityTooLarge         = HTTPStatus(http.StatusRequestEntityTooLarge)
	HTTPStatusRequestURITooLong             = HTTPStatus(http.StatusRequestURITooLong)
	HTTPStatusUnsupportedMediaType          = HTTPStatus(http.StatusUnsupportedMediaType)
	HTTPStatusRequestedRangeNotSatisfiable  = HTTPStatus(http.StatusRequestedRangeNotSatisfiable)
	HTTPStatusExpectationFailed             = HTTPStatus(http.StatusExpectationFailed)
	HTTPStatusTeapot                        = HTTPStatus(http.StatusTeapot)
	HTTPStatusUnprocessableEntity           = HTTPStatus(http.StatusUnprocessableEntity)
	HTTPStatusLocked                        = HTTPStatus(http.StatusLocked)
	HTTPStatusFailedDependency              = HTTPStatus(http.StatusFailedDependency)
	HTTPStatusUpgradeRequired               = HTTPStatus(http.StatusUpgradeRequired)
	HTTPStatusPreconditionRequired          = HTTPStatus(http.StatusPreconditionRequired)
	HTTPStatusTooManyRequests               = HTTPStatus(http.StatusTooManyRequests)
	HTTPStatusRequestHeaderFieldsTooLarge   = HTTPStatus(http.StatusRequestHeaderFieldsTooLarge)
	HTTPStatusUnavailableForLegalReasons    = HTTPStatus(http.StatusUnavailableForLegalReasons)
	HTTPStatusInternalServerError           = HTTPStatus(http.StatusInternalServerError)
	HTTPStatusNotImplemented                = HTTPStatus(http.StatusNotImplemented)
	HTTPStatusBadGateway                    = HTTPStatus(http.StatusBadGateway)
	HTTPStatusServiceUnavailable            = HTTPStatus(http.StatusServiceUnavailable)
	HTTPStatusGatewayTimeout                = HTTPStatus(http.StatusGatewayTimeout)
	HTTPStatusHTTPVersionNotSupported       = HTTPStatus(http.StatusHTTPVersionNotSupported)
	HTTPStatusVariantAlsoNegotiates         = HTTPStatus(http.StatusVariantAlsoNegotiates)
	HTTPStatusInsufficientStorage           = HTTPStatus(http.StatusInsufficientStorage)
	HTTPStatusLoopDetected                  = HTTPStatus(http.StatusLoopDetected)
	HTTPStatusNotExtended                   = HTTPStatus(http.StatusNotExtended)
	HTTPStatusNetworkAuthenticationRequired = HTTPStatus(http.StatusNetworkAuthenticationRequired)
)

// HTTPStatus returns a behavior that stores a HTTP status in the error metadata.
func HTTPStatus(status int) Behavior {
	return Metadata(reflect.ValueOf(HTTPStatus), status)
}

// GetHTTPStatus extracts a HTTP status from the error metadata, if any.
// It returns 0 if no HTTP status was set.
func GetHTTPStatus(err error) int {
	if status, ok := GetMetadata(err, reflect.ValueOf(HTTPStatus)).(int); ok {
		return status
	}
	return 0
}

// GetHTTPStatusOrDefault extracts a HTTP status from the error metadata, if any.
// It returns the given default HTTP status if no HTTP status was set.
func GetHTTPStatusOrDefault(err error, defaultStatus int) int {
	if status := GetHTTPStatus(err); status != 0 {
		return status
	}
	return defaultStatus
}
