package types

import (
	"context"
	"net/http"

	"github.com/ayushanand18/as-http3lib/internal/constants"
)

type ServeOptions struct {
	URL            string
	Handler        HandlerFunc
	Method         constants.HttpMethodTypes
	ResponseType   constants.ResponseTypes
	DefaultHeaders map[string]string
}

type HandlerFunc func(context.Context, *http.Request) (interface{}, error)
