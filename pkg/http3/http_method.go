package http3

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

type method struct {
	Method constants.HttpMethodTypes
	URL    string
	s      *server
}

type Method interface {
	Serve(types.ServeOptions)
}

func NewMethod(httpMethod constants.HttpMethodTypes, url string, s *server) Method {
	return &method{
		Method: httpMethod,
		URL:    url,
		s:      s,
	}
}

func (m *method) Serve(options types.ServeOptions) {
	if _, ok := m.s.routeMatchMap[m.URL]; !ok {
		m.s.routeMatchMap[m.URL] = make(map[constants.HttpMethodTypes]types.HandlerFunc)
	}

	if _, ok := m.s.routeMatchMap[m.URL]; !ok {
		m.s.routeMatchMap[m.URL] = make(map[constants.HttpMethodTypes]types.HandlerFunc)
	}

	// if the combination exists, reassign it
	m.s.routeMatchMap[m.URL][m.Method] = options.Handler
	if len(m.s.routeMatchMap[m.URL]) == 1 {
		m.s.mux.HandleFunc(m.URL, func(w http.ResponseWriter, r *http.Request) {
			requestMethod := constants.HttpMethodTypes(strings.ToUpper(r.Method))

			methodHandlers, ok := m.s.routeMatchMap[m.URL]
			if !ok {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			handler, ok := methodHandlers[requestMethod]
			if !ok {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			if options.Options.IsStreamingResponse {
				streamingDefaultHandler(r.Context(), w, handler, options.Decoder, options.Encoder, r)
			} else {
				httpDefaultHandler(r.Context(), w, handler, options.Decoder, options.Encoder, r)
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
