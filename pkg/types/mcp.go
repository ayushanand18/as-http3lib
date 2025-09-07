package types

type McpHandlerType int8

const (
	NoMcp     McpHandlerType = iota
	Tools     McpHandlerType = 1
	Resources McpHandlerType = 2
	Prompts   McpHandlerType = 3
)
