package http3

import (
	"context"
	"log/slog"

	"github.com/ayushanand18/as-http3lib/internal/mcp"
	"github.com/ayushanand18/as-http3lib/internal/utils"
	"github.com/ayushanand18/as-http3lib/pkg/errors"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

func (s *server) ListenAndServeMcp(ctx context.Context) error {
	utils.PrintStartBanner()
	errChan := make(chan error, 1)

	// Start the MCP server
	if err := s.registerAllMcpHandlers(ctx); err != nil {
		return err
	}

	go func() {
		slog.InfoContext(ctx, "Starting MCP Server", "port", s.mcpServer.Addr)
		// server over https
		errChan <- s.mcpServer.ListenAndServeTLS("", "")
	}()

	return <-errChan
}

func (s *server) registerAllMcpHandlers(ctx context.Context) error {
	// register on the /mcp route
	s.POST("/mcp").Serve(types.ServeOptions{
		Handler: s.mcpRootHandler(),
	})

	return nil
}

// MCP_TOOL takes :path: as argument, and registers a tool method
func (s *server) MCP_TOOL(path string) Method {
	return nil
}

func (s *server) MCP_RESOURCE(path string) Method {
	return nil
}

func (s *server) MCP_PROMPT(path string) Method {
	return nil
}

func (s *server) mcpRootHandler() types.HandlerFunc {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		// version checks
		if err := validateJsonRpcVersion(ctx, req); err != nil {
			return nil, err
		}

		// id and method check
		return handleMcpRequest(ctx, req)
	}
}

func validateJsonRpcVersion(ctx context.Context, req interface{}) error {
	return nil
}

func handleMcpRequest(ctx context.Context, req interface{}) (interface{}, error) {
	request, ok := req.(mcp.McpRequest)
	if !ok {
		return nil, errors.InternalServerError.New("Error while parsing MCP Request")
	}

	handler, ok := mcpMethodToHandlerMap[request.Method]
	if !ok {
		return defaultUnhandledHandler(ctx, request)
	}

	return handler(ctx, request)
}

func handleMcpInitialize(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	return mcp.McpResponse{}, nil
}

func defaultUnhandledHandler(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	return mcp.McpResponse{}, errors.MethodNotAllowed.New("Method not allowed")
}
