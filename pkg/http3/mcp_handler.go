package http3

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ayushanand18/as-http3lib/internal/constants"
	"github.com/ayushanand18/as-http3lib/internal/mcp"
	"github.com/ayushanand18/as-http3lib/pkg/errors"
	"github.com/ayushanand18/as-http3lib/pkg/types"
)

func (s *server) defaultUnhandledHandler(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	return s.handleMcpInitialize(ctx, req)
}

func (s *server) handleMcpInitialize(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	subscribeTrue := true
	data := mcp.McpResponse{
		JsonRpc: "2.0",
		Id:      0,
		Result: mcp.InitializeResult{
			ProtocolVersion: "2025-06-18",
			Capabilities: mcp.McpCapabilities{
				Tools: mcp.McpTools{
					ListChanged: true,
				},
				Resources: mcp.McpResources{
					Subscribe:   &subscribeTrue,
					ListChanged: true,
				},
				Prompts: mcp.McpPrompts{
					ListChanged: true,
				},
				Logging: map[string]interface{}{},
			},
			ServerInfo: mcp.McpServerInfo{
				Name:    "Ash3 - Mcp Server",
				Version: "2.0.0",
				Title:   "ashttp3-lib Mcp Server",
			},
		},
	}

	dataStr, err := json.Marshal(data)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal mcp initialize response", "error", err)
		return mcp.McpResponse{}, errors.InternalServerError.New("failed to marshal mcp initialize response: " + err.Error())
	}

	channel, ok := ctx.Value(constants.StreamingResponseChannelContextKey).(chan types.StreamChunk)
	if ok {
		channel <- types.StreamChunk{
			Id:   uint32(1),
			Data: dataStr,
		}
	}

	return data, nil
}

func (s *server) handleMcpNotificationsInitialized(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	return mcp.McpResponse{}, nil
}

func (s *server) handleMcpToolsList(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	tools := make([]mcp.McpTool, 0, len(s.mcpTools))
	for _, tool := range s.mcpTools {
		tools = append(tools, mcp.McpTool{
			Name:         tool.Name,
			Description:  tool.Description,
			InputSchema:  tool.InputSchema,
			OutputSchema: tool.OutputSchema,
		})
	}

	return mcp.McpResponse{
		JsonRpc: "2.0",
		Id:      req.Id,
		Result: mcp.McpToolsListResult{
			Tools: tools,
		},
	}, nil
}

func (s *server) handleMcpToolsCall(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	var params mcp.McpToolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return mcp.McpResponse{}, errors.BadRequest.New("invalid tools/call params: " + err.Error())
	}

	tool, ok := s.mcpTools[params.Name]
	if !ok {
		return mcp.McpResponse{}, errors.NotFound.New("tool not found: " + params.Name)
	}

	result, err := tool.Handler(ctx, params.Arguments)
	if err != nil {
		return mcp.McpResponse{}, err
	}

	return mcp.McpResponse{
		JsonRpc: "2.0",
		Id:      req.Id,
		Result: mcp.McpToolCallResult{
			Content: result,
		},
	}, nil
}

func (s *server) handleMcpResourcesList(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	resources := make([]mcp.McpResource, 0, len(s.mcpResources))
	for _, r := range s.mcpResources {
		resources = append(resources, r)
	}

	return mcp.McpResponse{
		JsonRpc: "2.0",
		Id:      req.Id,
		Result: map[string]interface{}{
			"resources": resources,
		},
	}, nil
}

func (s *server) handleMcpResourcesRead(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	var params mcp.McpResourceReadParams
	if _, err := CastParams(req.Params, &params); err != nil {
		return mcp.McpResponse{}, errors.BadRequest.New("invalid resources/read params: " + err.Error())
	}

	resource, ok := s.mcpResources[params.URI]
	if !ok {
		return mcp.McpResponse{}, errors.NotFound.New("resource not found: " + params.URI)
	}

	result, err := resource.Handler(ctx, nil)
	if err != nil {
		return mcp.McpResponse{}, err
	}

	return mcp.McpResponse{
		JsonRpc: "2.0",
		Id:      req.Id,
		Result: map[string]interface{}{
			"uri":      resource.URI,
			"name":     resource.Name,
			"mimeType": resource.MimeType,
			"content":  result,
		},
	}, nil
}

func (s *server) handleMcpPromptsList(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	prompts := make([]mcp.McpPrompt, 0, len(s.mcpPrompts))
	for _, p := range s.mcpPrompts {
		prompts = append(prompts, p)
	}

	return mcp.McpResponse{
		JsonRpc: "2.0",
		Id:      req.Id,
		Result: map[string]interface{}{
			"listChanged": false,
			"prompts":     prompts,
		},
	}, nil
}

func (s *server) handleMcpPromptsGet(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	var params mcp.McpPromptGetParams
	_, err := CastParams(req.Params, &params)
	if err != nil {
		return mcp.McpResponse{}, errors.BadRequest.New("invalid prompts/get params: " + err.Error())
	}

	prompt, ok := s.mcpPrompts[params.Name]
	if !ok {
		return mcp.McpResponse{}, errors.NotFound.New("prompt not found: " + params.Name)
	}

	return mcp.McpResponse{
		JsonRpc: "2.0",
		Id:      req.Id,
		Result: map[string]interface{}{
			"prompt": prompt,
		},
	}, nil
}

func (s *server) handleMcpShutdown(ctx context.Context, req mcp.McpRequest) (mcp.McpResponse, error) {
	slog.InfoContext(ctx, "Received MCP shutdown request", "req", req)

	return mcp.McpResponse{
		JsonRpc: "2.0",
		Id:      req.Id,
		Result: map[string]interface{}{
			"message": "Server shutdown initiated",
		},
	}, nil
}

func healthCheckHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return map[string]string{
		"status": "ok",
	}, nil
}
