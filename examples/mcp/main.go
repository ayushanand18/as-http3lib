package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ayushanand18/crazyhttp/pkg/errors"
	"github.com/ayushanand18/crazyhttp/pkg/http3"
	"github.com/ayushanand18/crazyhttp/pkg/types"
)

func resourceHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return map[string]interface{}{
		"uri":     "my_file",
		"content": "This is the file content",
	}, nil
}

func promptHandler(ctx context.Context, req interface{}) (interface{}, error) {
	args, ok := req.(map[string]interface{})
	if !ok {
		return nil, errors.BadRequest.New("invalid args for prompt")
	}

	name, _ := args["name"].(string)
	if name == "" {
		name = "Guest"
	}

	return map[string]interface{}{
		"message": fmt.Sprintf("Hello, %s! Welcome to MCP server.", name),
	}, nil
}

var promptArgsSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": "Name of the user to greet",
		},
	},
	"required": []string{"name"},
}

func main() {
	ctx := context.Background()

	server := http3.NewServer(ctx)
	if err := server.Initialize(ctx); err != nil {
		log.Fatalf("Server failed to Initialize: %v", err)
	}

	server.MCP_TOOL("echo").
		WithDescription("Echo back input").
		WithInputSchema(map[string]string{"message": "string"}).
		Serve(types.ServeOptions{
			Handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				args := req.(map[string]interface{})
				return map[string]string{"echo": args["message"].(string)}, nil
			},
		})

	server.MCP_RESOURCE("my_file").
		WithName("file-reader").
		WithDescription("Reads contents of a file").
		WithMimeType("text/plain").
		Serve(types.ServeOptions{Handler: resourceHandler})

	server.MCP_PROMPT("greeting").
		WithDescription("Returns greeting messages").
		WithInputSchema(promptArgsSchema).
		Serve(types.ServeOptions{Handler: promptHandler})

	if err := server.ListenAndServeMcp(ctx); err != nil {
		log.Fatalf("MCP Server failed to start: %v", err)
	}
}
