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
	HTTPStatusMisdirectedRequest            = HTTPStatus(http.StatusMisdirectedRequest)
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

var httpPublicMessages = map[int]Behavior{
	http.StatusBadRequest:                    PublicMessage("bad-request"),
	http.StatusUnauthorized:                  PublicMessage("unauthorized"),
	http.StatusPaymentRequired:               PublicMessage("payment-required"),
	http.StatusForbidden:                     PublicMessage("forbidden"),
	http.StatusNotFound:                      PublicMessage("not-found"),
	http.StatusMethodNotAllowed:              PublicMessage("method-not-allowed"),
	http.StatusNotAcceptable:                 PublicMessage("not-acceptable"),
	http.StatusProxyAuthRequired:             PublicMessage("proxy-auth-required"),
	http.StatusRequestTimeout:                PublicMessage("request-timeout"),
	http.StatusConflict:                      PublicMessage("conflict"),
	http.StatusGone:                          PublicMessage("gone"),
	http.StatusLengthRequired:                PublicMessage("length-required"),
	http.StatusPreconditionFailed:            PublicMessage("precondition-failed"),
	http.StatusRequestEntityTooLarge:         PublicMessage("request-entity-too-large"),
	http.StatusRequestURITooLong:             PublicMessage("request-uri-too-long"),
	http.StatusUnsupportedMediaType:          PublicMessage("unsupported-media-type"),
	http.StatusRequestedRangeNotSatisfiable:  PublicMessage("requested-range-not-satisfiable"),
	http.StatusExpectationFailed:             PublicMessage("expectation-failed"),
	http.StatusTeapot:                        PublicMessage("i-am-a-teapot"),
	http.StatusMisdirectedRequest:            PublicMessage("misdirected-request"),
	http.StatusUnprocessableEntity:           PublicMessage("unprocessable-entity"),
	http.StatusLocked:                        PublicMessage("locked"),
	http.StatusFailedDependency:              PublicMessage("failed-dependency"),
	http.StatusUpgradeRequired:               PublicMessage("upgrade-required"),
	http.StatusPreconditionRequired:          PublicMessage("precondition-required"),
	http.StatusTooManyRequests:               PublicMessage("too-many-requests"),
	http.StatusRequestHeaderFieldsTooLarge:   PublicMessage("request-header-fields-too-large"),
	http.StatusUnavailableForLegalReasons:    PublicMessage("unavailable-for-legal-reasons"),
	http.StatusInternalServerError:           PublicMessage("internal-server-error"),
	http.StatusNotImplemented:                PublicMessage("not-implemented"),
	http.StatusBadGateway:                    PublicMessage("bad-gateway"),
	http.StatusServiceUnavailable:            PublicMessage("service-unavailable"),
	http.StatusGatewayTimeout:                PublicMessage("gateway-timeout"),
	http.StatusHTTPVersionNotSupported:       PublicMessage("http-version-not-supported"),
	http.StatusVariantAlsoNegotiates:         PublicMessage("variant-also-negotiates"),
	http.StatusInsufficientStorage:           PublicMessage("insufficient-storage"),
	http.StatusLoopDetected:                  PublicMessage("loop-detected"),
	http.StatusNotExtended:                   PublicMessage("not-extended"),
	http.StatusNetworkAuthenticationRequired: PublicMessage("network-authentication-required"),
}

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

// HTTPPublicMessage returns a default PublicMessage Behavior corresponding to the given HTTP status.
// If the given status is not a HTTP 4xx or 5xx status registered with IANA, it returns PublicMessage("unknown").
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
func HTTPPublicMessage(status int) Behavior {
	if behavior, ok := httpPublicMessages[status]; ok {
		return behavior
	}
	return PublicMessage("unknown")
}

// HTTPError returns a compound Behavior that includes both HTTPStatus and HTTPublicMessage for the given HTTP status.
func HTTPError(status int) Behavior {
	return Behaviors(HTTPStatus(status), HTTPPublicMessage(status))
}
