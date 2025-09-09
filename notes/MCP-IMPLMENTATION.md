## MCP Implementation and Architecture

1. tools
2. resources
3. prompts

code 
```golang

server.MCP_TOOL("/tool").Serve(types.ServeOptions{
    Handler: mcpToolsHandler,
    Encoder: Encoder,
})
```

* default endpoints that need to be populated
path                | description
--------------------|--------------------------
tools/list          | list all tools
resources/list      | list all resources
prompts/list        | list all prompts

