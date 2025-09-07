package mcp

import (
	"context"
	"encoding/json"

	"github.com/ayushanand18/as-http3lib/pkg/types"
)

type McpRequest struct {
	JsonRpc string          `json:"jsonrpc"`
	Id      string          `json:"id"`
	Method  McpMethodTypes  `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type McpParams struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    McpCapabilities `json:"capabilities"`
	ClientInfo      McpClientInfo   `json:"clientInfo"`
}

type McpCapabilities struct {
	Experimental bool                   `json:"experimental"`
	Tools        McpTools               `json:"tools"`
	Resources    McpResources           `json:"resources"`
	Prompts      McpPrompts             `json:"prompts"`
	Logging      map[string]interface{} `json:"logging"`
}

type McpClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type McpResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *McpError   `json:"error,omitempty"`
}

type McpError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type McpMethodTypes string

const (
	Initialize               McpMethodTypes = "initialize"
	NotificationsInitialized McpMethodTypes = "notifications/initialized"
	ToolsList                McpMethodTypes = "tool/list"
	ToolsCall                McpMethodTypes = "tool/call"
	ResourcesList            McpMethodTypes = "resources/list"
	ResourcesRead            McpMethodTypes = "resources/read"
	PromptsList              McpMethodTypes = "prompts/list"
	PromptsGet               McpMethodTypes = "prompts/get"
	Shutdown                 McpMethodTypes = "shutdown"
)

type InitializeResult struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    McpCapabilities `json:"capabilities"`
	ServerInfo      McpServerInfo   `json:"serverInfo"`
}

type McpServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Title   string `json:"title"`
}

type McpTools struct {
	ListChanged bool  `json:"listChanged"`
	Subscribe   *bool `json:"subscribe,omitempty"`
}
type McpResources struct {
	ListChanged bool  `json:"listChanged"`
	Subscribe   *bool `json:"subscribe,omitempty"`
}

type McpPrompts struct {
	ListChanged bool  `json:"listChanged"`
	Subscribe   *bool `json:"subscribe,omitempty"`
}

type McpTool struct {
	Name         string
	Description  string
	InputSchema  interface{}
	OutputSchema interface{}
	Handler      types.HandlerFunc
}

type McpResource struct {
	URI         string
	Name        string
	MimeType    string
	Description string
	Handler     types.HandlerFunc
}

type McpPrompt struct {
	Name        string
	Description string
	Messages    []Message
	ArgsSchema  interface{}
}

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type McpPromptGetParams struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args,omitempty"`
}

type McpHandlerType int8

const (
	NoMcp     McpHandlerType = iota
	Tools     McpHandlerType = 1
	Resources McpHandlerType = 2
	Prompts   McpHandlerType = 3
)

type McpHandlerFunc func(context.Context, McpRequest) (McpResponse, error)

type McpInitializeParams struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    McpCapabilities `json:"capabilities"`
	ClientInfo      McpClientInfo   `json:"clientInfo"`
}

type McpToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type McpResourceReadParams struct {
	URI string `json:"uri"`
}

type McpToolCallResult struct {
	Content interface{} `json:"content"`
}

type McpToolsListResult struct {
	Tools []McpTool `json:"tools"`
}
