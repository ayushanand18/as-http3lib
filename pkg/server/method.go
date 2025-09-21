package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ayushanand18/crazyhttp/internal/constants"
	"github.com/ayushanand18/crazyhttp/internal/ratelimiter"
	"github.com/ayushanand18/crazyhttp/pkg/types"
)

type method struct {
	Method constants.HttpMethodTypes
	URL    string
	s      *server

	// utility
	rateLimiter *ratelimiter.RateLimiter

	description  string
	inputSchema  interface{}
	outputSchema interface{}
	name         string
	handler      types.HandlerFunc
}

type Method interface {
	Serve(types.ServeOptions)

	WithRateLimit(types.RateLimitOptions) Method
	WithDescription(desc string) Method
	WithInputSchema(schema interface{}) Method
	WithOutputSchema(schema interface{}) Method
	WithName(name string) Method
}

func NewMethod(httpMethod constants.HttpMethodTypes, url string, s *server) Method {
	return &method{
		Method: httpMethod,
		URL:    url,
		s:      s,
	}
}

func (m *method) Serve(options types.ServeOptions) {
	m.handler = options.Handler

	if _, ok := m.s.routeMatchMap[m.URL]; !ok {
		m.s.routeMatchMap[m.URL] = make(map[constants.HttpMethodTypes]types.ServeOptions)
	}

	if _, ok := m.s.routeMatchMap[m.URL]; !ok {
		m.s.routeMatchMap[m.URL] = make(map[constants.HttpMethodTypes]types.ServeOptions)
	}

	// if the combination exists, reassign it
	m.s.routeMatchMap[m.URL][m.Method] = options
	if len(m.s.routeMatchMap[m.URL]) == 1 {
		m.s.mux.HandleFunc(m.URL, func(w http.ResponseWriter, r *http.Request) {
			requestMethod := constants.HttpMethodTypes(strings.ToUpper(r.Method))

			methodHandlers, ok := m.s.routeMatchMap[m.URL]
			if !ok {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			opts, ok := methodHandlers[requestMethod]
			if !ok {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			DumpRequest(r)
			if options.Options.IsStreamingResponse {
				streamingDefaultHandler(r.Context(), w, opts.Handler, opts.Decoder, opts.Encoder, r, m)
			} else {
				httpDefaultHandler(r.Context(), w, opts.Handler, opts.Decoder, opts.Encoder, r, m)
			}
		})
	}
}

func DecodeJsonRequest[T any](in interface{}) (T, error) {
	var out T
	raw, err := json.Marshal(in)
	if err != nil {
		return out, err
	}

	err = json.Unmarshal(raw, &out)
	return out, err
}

func (m *method) WithDescription(desc string) Method {
	m.description = desc
	return m
}

func (m *method) WithInputSchema(schema interface{}) Method {
	m.inputSchema = schema
	return m
}

func (m *method) WithOutputSchema(schema interface{}) Method {
	m.outputSchema = schema
	return m
}

func (m *method) WithName(name string) Method {
	m.name = name
	return m
}

func (m *method) WithRateLimit(options types.RateLimitOptions) Method {
	m.rateLimiter = ratelimiter.NewRateLimiter(options.Limit, time.Duration(options.BucketDurationInSeconds)*time.Second)

	return m
}
