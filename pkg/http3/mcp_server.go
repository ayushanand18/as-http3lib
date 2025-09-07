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
		// server over http
		errChan <- s.mcpServer.ListenAndServe()
	}()

	return <-errChan
}

func (s *server) registerAllMcpHandlers(ctx context.Context) error {
	s.mcpMethodToHandlerMap = map[mcp.McpMethodTypes]mcp.McpHandlerFunc{
		mcp.Initialize:               s.handleMcpInitialize,
		mcp.NotificationsInitialized: s.handleMcpNotificationsInitialized,
		mcp.ToolsList:                s.handleMcpToolsList,
		mcp.ToolsCall:                s.handleMcpToolsCall,
		mcp.ResourcesList:            s.handleMcpResourcesList,
		mcp.ResourcesRead:            s.handleMcpResourcesRead,
		mcp.PromptsList:              s.handleMcpPromptsList,
		mcp.PromptsGet:               s.handleMcpPromptsGet,
		mcp.Shutdown:                 s.handleMcpShutdown,
	}

	s.POST("/mcp").Serve(types.ServeOptions{
		Handler: s.mcpRootHandler(),
	})

	return nil
}

func (s *server) mcpRootHandler() types.HandlerFunc {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		// version checks
		if err := validateJsonRpcVersion(ctx, req); err != nil {
			return nil, err
		}

		// id and method check
		return s.handleMcpRequest(ctx, req)
	}
}

func validateJsonRpcVersion(ctx context.Context, req interface{}) error {
	return nil
}

func (s *server) handleMcpRequest(ctx context.Context, req interface{}) (interface{}, error) {
	request, err := DecodeJsonRequest[mcp.McpRequest](req)
	if err != nil {
		return nil, errors.InternalServerError.New("Error while parsing MCP Request")
	}

	handler, ok := s.mcpMethodToHandlerMap[request.Method]
	if !ok {
		return s.defaultUnhandledHandler(ctx, request)
	}

	return handler(ctx, request)
}

// MCP_TOOL takes :path: as argument, and registers a tool method
func (s *server) MCP_TOOL(path string) Method {
	if s.mcpTools == nil {
		s.mcpTools = make(map[string]mcp.McpTool)
	}
	return &method{
		URL:     path,
		s:       s,
		mcpKind: "tool",
	}
}

func (s *server) MCP_RESOURCE(path string) Method {
	return &method{
		URL:     path,
		s:       s,
		mcpKind: "resource",
	}
}

func (s *server) MCP_PROMPT(path string) Method {
	return &method{
		URL:     path,
		s:       s,
		mcpKind: "prompt",
	}
}

func (s *server) ExecuteTool(ctx context.Context, name string, req interface{}) (interface{}, error) {
	tool, ok := s.mcpTools[name]
	if !ok {
		return nil, errors.NotFound.New("tool not found: " + name)
	}
	return tool.Handler(ctx, req)
}
