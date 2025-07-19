package constants

type HttpMethodTypes string

const (
	HTTP_METHOD_GET     HttpMethodTypes = "GET"
	HTTP_METHOD_POST    HttpMethodTypes = "POST"
	HTTP_METHOD_PUT     HttpMethodTypes = "PUT"
	HTTP_METHOD_PATCH   HttpMethodTypes = "PATCH"
	HTTP_METHOD_DELETE  HttpMethodTypes = "DELETE"
	HTTP_METHOD_HEAD    HttpMethodTypes = "HEAD"
	HTTP_METHOD_OPTIONS HttpMethodTypes = "OPTIONS"
	HTTP_METHOD_CONNECT HttpMethodTypes = "CONNECT"
	HTTP_METHOD_TRACE   HttpMethodTypes = "TRACE"
)

type ResponseTypes int

const (
	RESPONSE_TYPE_BASE_RESPONSE      ResponseTypes = iota
	RESPONSE_TYPE_STREAMING_RESPONSE ResponseTypes = 1
	RESPONSE_TYPE_JSON_RESPONSE      ResponseTypes = 2
)

type ContextKeys string

const (
	STREAMING_RESPONSE_CHANNEL_CONTEXT_KEY ContextKeys = "response_channel"
	HTTP_REQUEST_HEADERS                   ContextKeys = "request_headers"
	HTTP_REQUEST_URL_PARAMS                ContextKeys = "request_url_params"
	HTTP_REQUEST_PATH_VALUES               ContextKeys = "request_path_values"
)
