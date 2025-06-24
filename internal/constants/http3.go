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
