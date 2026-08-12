package handlers

import (
	"fmt"

	"github.com/chilly-ai-pilot/go-search-mcp/config"
	"github.com/mark3labs/mcp-go/server"
)

type HandlerBuilder func(cfg *config.Config) (server.ToolHandlerFunc, error)

var handlerBuilders = map[string]HandlerBuilder{}

func RegisterHandler(handlerType string, builder HandlerBuilder) {
	handlerBuilders[handlerType] = builder
}

type Registry struct {
	handlers map[string]server.ToolHandlerFunc
}

func NewRegistry(cfg *config.Config, tools []config.ToolDefinition) (*Registry, error) {
	handlers := make(map[string]server.ToolHandlerFunc, len(tools))
	for _, tool := range tools {
		builder, ok := handlerBuilders[tool.Name]
		if !ok {
			return nil, fmt.Errorf("unknown handler type %q", tool.Name)
		}
		handler, err := builder(cfg)
		if err != nil {
			return nil, fmt.Errorf("initialize handler %q: %w", tool.Name, err)
		}
		handlers[tool.Name] = handler
	}

	return &Registry{handlers: handlers}, nil
}

func (r *Registry) Get(name string) (server.ToolHandlerFunc, error) {
	handler, ok := r.handlers[name]
	if !ok {
		return nil, fmt.Errorf("unknown handler: %s", name)
	}
	return handler, nil
}
