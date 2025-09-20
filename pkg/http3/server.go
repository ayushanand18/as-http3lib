package http3

import (
	"context"
	"net/http"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/internal/mcp"
	"github.com/ayushanand18/as-http3lib/internal/utils"
	"github.com/ayushanand18/as-http3lib/pkg/types"
	"github.com/gorilla/mux"
	"github.com/quic-go/quic-go"
	qchttp3 "github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
)

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range injectConstantHeaders() {
		w.Header().Set(k, v)
	}

	rec := &responseRecorder{ResponseWriter: w, status: 0}
	h.mux.ServeHTTP(rec, r)

	if !rec.wrote {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("404 page not found\n"))
	}
}

type server struct {
	// hTTP server assets
	h3server       qchttp3.Server
	mux            *mux.Router
	routeMatchMap  map[string]map[constants.HttpMethodTypes]types.ServeOptions
	http1ServerTLS http.Server
	http1Server    http.Server

	// MCP server assets
	mcpServer             http.Server
	mcpMethodToHandlerMap map[mcp.McpMethodTypes]mcp.McpHandlerFunc
	mcpTools              map[string]mcp.McpTool
	mcpResources          map[string]mcp.McpResource
	mcpPrompts            map[string]mcp.McpPrompt
}

type Server interface {
	Initialize(context.Context) error
	ListenAndServe(context.Context) error
	ListenAndServeMcp(context.Context) error

	// HTTP Methods
	GET(string) Method
	POST(string) Method
	PUT(string) Method
	PATCH(string) Method
	DELETE(string) Method
	HEAD(string) Method
	OPTIONS(string) Method
	CONNECT(string) Method
	TRACE(string) Method

	// MCP Methods
	MCP_TOOL(string) Method
	MCP_RESOURCE(string) Method
	MCP_PROMPT(string) Method
}

func NewServer(ctx context.Context) Server {
	quicConfig := &quic.Config{
		Tracer:          qlog.DefaultConnectionTracer,
		Allow0RTT:       true,
		EnableDatagrams: true,
	}
	return &server{
		h3server: qchttp3.Server{
			Addr:            utils.GetListeningAddress(ctx),
			Handler:         nil,
			EnableDatagrams: true,
			QUICConfig:      quicConfig,
		},
		http1Server: http.Server{
			Addr: utils.GetHttp1ListeningAddress(ctx),
		},
		http1ServerTLS: http.Server{
			Addr: utils.GetHttp1TLSListeningAddress(ctx),
		},
		mcpServer: http.Server{
			Addr: utils.GetMcpListeningAddress(ctx),
		},
		mux:           mux.NewRouter(),
		routeMatchMap: make(map[string]map[constants.HttpMethodTypes]types.ServeOptions),
	}
}
