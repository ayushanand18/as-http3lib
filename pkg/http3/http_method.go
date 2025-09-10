package http3

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/internal/mcp"
	"github.com/ayushanand18/as-http3lib/internal/ratelimiter"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

type method struct {
	Method constants.HttpMethodTypes
	URL    string
	s      *server

	// utility
	rateLimiter *ratelimiter.RateLimiter

	// MCP specific fields (optional, only used for MCP_TOOL/RESOURCE/PROMPT)
	mcpKind      string // "tool" | "resource" | "prompt"
	description  string
	inputSchema  interface{}
	outputSchema interface{}
	mimeType     string // only for resource
	name         string
	handler      types.HandlerFunc
}

type Method interface {
	Serve(types.ServeOptions)
	WithRateLimit(types.RateLimitOptions) Method
	WithDescription(desc string) Method
	WithInputSchema(schema interface{}) Method
	WithOutputSchema(schema interface{}) Method
	WithMimeType(mime string) Method
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

	switch m.mcpKind {
	case "tool":
		m.s.mcpTools[m.URL] = mcp.McpTool{
			Name:         m.URL,
			Description:  m.description,
			InputSchema:  m.inputSchema,
			OutputSchema: m.outputSchema,
			Handler:      m.handler,
		}
	case "resource":
		if m.s.mcpResources == nil {
			m.s.mcpResources = make(map[string]mcp.McpResource)
		}
		m.s.mcpResources[m.URL] = mcp.McpResource{
			URI:         m.URL,
			Name:        m.name,
			MimeType:    m.mimeType,
			Description: m.description,
			Handler:     m.handler,
		}

	case "prompt":
		if m.s.mcpPrompts == nil {
			m.s.mcpPrompts = make(map[string]mcp.McpPrompt)
		}
		m.s.mcpPrompts[m.URL] = mcp.McpPrompt{
			Name:        m.name,
			Description: m.description,
			ArgsSchema:  m.inputSchema,
		}
	default:
		// fall back to HTTP registration
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
					streamingDefaultHandler(r.Context(), w, handler, options.Decoder, options.Encoder, r, m)
				} else {
					httpDefaultHandler(r.Context(), w, handler, options.Decoder, options.Encoder, r, m)
				}
			})
		}
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

func (m *method) WithMimeType(mime string) Method {
	m.mimeType = mime
	return m
}

func (m *method) WithRateLimit(options types.RateLimitOptions) Method {
	m.rateLimiter = ratelimiter.NewRateLimiter(options.Limit, time.Duration(options.BucketDurationInSeconds)*time.Second)

	return m
}
