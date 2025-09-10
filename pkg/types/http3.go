package types

import (
	"context"
	"net/http"
)

type ServeOptions struct {
	Handler     HandlerFunc
	Decoder     HttpDecoder
	Encoder     HttpEncoder
	BeforeServe HttpRequestMiddleware
	AfterServer HttpResponseMiddleware
	Options     MethodOptions
}

type MethodOptions struct {
	IsStreamingResponse bool
}

type HandlerFunc func(context.Context, interface{}) (interface{}, error)

type HttpDecoder func(ctx context.Context, r *http.Request) (request interface{}, err error)

type HttpEncoder func(ctx context.Context, response interface{}) (headers map[string][]string, body []byte, err error)

type HttpRequestMiddleware func(ctx context.Context, incomingRequest interface{}) (outgoingContext context.Context, outgoingRequest interface{}, err error)

type HttpResponseMiddleware func(ctx context.Context, incomingResponse interface{}) (outgoingResponse interface{}, err error)

type RateLimitOptions struct {
	Limit                   int    // number of requests allowed in the given duration
	BucketDurationInSeconds int64  // duration in seconds for which the limit is applicable
	ContextKey              string // key in context which will be checked for rate limiting
}
