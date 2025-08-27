package http3

import "context"

func (s *server) ListenAndServeMcp(ctx context.Context) error {
	// Start the MCP server
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
