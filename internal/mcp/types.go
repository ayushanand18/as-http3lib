package mcp

type McpRequest struct {
	JsonRpc string         `json:"jsonrpc"`
	Id      string         `json:"id"`
	Method  McpMethodTypes `json:"method"`
	Params  McpParams      `json:"params"`
}

type McpParams struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    McpCapabilities `json:"capabilities"`
	ClientInfo      McpClientInfo   `json:"clientInfo"`
}

type McpCapabilities struct {
	Experimental bool `json:"experimental"`
}

type McpClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type McpResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result"`
}

type McpMethodTypes string

const (
	Initialize               McpMethodTypes = "initialize"
	NotificationsInitialized McpMethodTypes = "notifications/initialized"
	ToolsList                McpMethodTypes = "tools/list"
	ToolsCall                McpMethodTypes = "tools/call"
	ResourcesList            McpMethodTypes = "resources/list"
	ResourcesRead            McpMethodTypes = "resources/read"
	PromptsList              McpMethodTypes = "prompts/list"
	PromptsGet               McpMethodTypes = "prompts/get"
	Shutdown                 McpMethodTypes = "shutdown"
)
