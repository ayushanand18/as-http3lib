package errors

import "net/http"

type ErrorType uint32

const (
	NoType ErrorType = iota

	// 4xx
	BadRequest
	Unauthorized
	PaymentRequired
	Forbidden
	NotFound
	MethodNotAllowed
	NotAcceptable
	ProxyAuthenticationRequired
	RequestTimeout
	Conflict
	LengthRequired
	PreconditionFailed
	ContentTooLarge
	URITooLong
	UnsupportedMediaType
	RangeNotSatisfiable
	FailedDependency
	TooEarly
	TooManyRequests
	UnavailableForLegalReasons

	// 5xx
	InternalServerError
	NotImplemented
	BadGateway
	ServiceUnavailable
	GatewayTimeout
	HTTPVersionNotSupported
	VariantAlsoNegotiates
	InsufficientStorage
	LoopDetected
	NotExtended
	NetworkAuthenticationRequired
)

var errorTypeToMessageMap = map[ErrorType]string{
	BadRequest:                  "BAD_REQUEST_ERROR",
	Unauthorized:                "UNAUTHORISED_ERROR",
	PaymentRequired:             "PAYMENT_REQUIRED_ERROR",
	Forbidden:                   "FORBIDDEN_ERROR",
	NotFound:                    "NOT_FOUND_ERROR",
	MethodNotAllowed:            "METHOD_NOT_ALLOWED_ERROR",
	NotAcceptable:               "NOT_ACCEPTABLE_ERROR",
	ProxyAuthenticationRequired: "PROXY_AUTHENTICATION_REQUIRED_ERROR",
	RequestTimeout:              "REQUEST_ERROR",
	Conflict:                    "CONFLICT_ERROR",
	LengthRequired:              "LENGTH_REQUIRED_ERROR",
	PreconditionFailed:          "PRECONDITION_FAILED_ERROR",
	ContentTooLarge:             "CONTENT_TOO_LARGE_ERROR",
	URITooLong:                  "URI_TOO_LONG_ERROR",
	UnsupportedMediaType:        "UNSUPPORTED_MEDIA_TYPE_ERROR",
	RangeNotSatisfiable:         "RANGE_NOT_SATISFIABLE_ERROR",
	FailedDependency:            "FAILED_DEPENDENCY_ERROR",
	TooEarly:                    "TOO_EARLY_ERROR",
	TooManyRequests:             "TOO_MANY_REQUESTS_ERROR",
	UnavailableForLegalReasons:  "UNAVAILABLE_FOR_LEGAL_REASONS_ERROR",

	InternalServerError:           "INTERNAL_SERVER_ERROR",
	NotImplemented:                "NOT_IMPLEMENTED_ERROR",
	BadGateway:                    "BAD_GATEWAY_ERROR",
	ServiceUnavailable:            "SERVICE_UNAVAILABLE_ERROR",
	GatewayTimeout:                "GATEWAY_TIMEOUT_ERROR",
	HTTPVersionNotSupported:       "HTTP_VERSION_NOT_SUPPORTED_ERROR",
	VariantAlsoNegotiates:         "VARIANT_ALSO_NEGOTIATES_ERROR",
	InsufficientStorage:           "INSUFFICIENT_STORAGE_ERROR",
	LoopDetected:                  "LOOP_DETECTED_ERROR",
	NotExtended:                   "NOT_EXTENDED_ERROR",
	NetworkAuthenticationRequired: "NETWORK_AUTHENTICATION_REQUIRED_ERROR",
}

var errorTypeToStatusCodeMap = map[ErrorType]int{
	BadRequest:                  400,
	Unauthorized:                401,
	PaymentRequired:             402,
	Forbidden:                   403,
	NotFound:                    404,
	MethodNotAllowed:            405,
	NotAcceptable:               406,
	ProxyAuthenticationRequired: 407,
	RequestTimeout:              408,
	Conflict:                    409,
	LengthRequired:              411,
	PreconditionFailed:          412,
	ContentTooLarge:             413,
	URITooLong:                  414,
	UnsupportedMediaType:        415,
	RangeNotSatisfiable:         416,
	FailedDependency:            424,
	TooEarly:                    425,
	TooManyRequests:             429,
	UnavailableForLegalReasons:  451,

	InternalServerError:           500,
	NotImplemented:                501,
	BadGateway:                    502,
	ServiceUnavailable:            503,
	GatewayTimeout:                504,
	HTTPVersionNotSupported:       505,
	VariantAlsoNegotiates:         506,
	InsufficientStorage:           507,
	LoopDetected:                  508,
	NotExtended:                   510,
	NetworkAuthenticationRequired: 511,
}

func DecodeErrorToHttpErrorStatus(err error) int {
	errTyped, ok := err.(customError)
	if !ok {
		return http.StatusInternalServerError
	}

	return errTyped.Code()
}
