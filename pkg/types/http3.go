package types

import (
	"context"
	"net/http"

	"github.com/ayushanand18/as-http3lib/internal/constants"
)

type ServeOptions struct {
	URL     string
	Handler HandlerFunc
	Method  constants.HttpMethodTypes
}

type HttpResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

type HandlerFunc func(context.Context, *http.Request) *HttpResponse
