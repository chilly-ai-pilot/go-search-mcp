package config

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

type ToolDefinition struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	InputSchema ToolInputSchemaConfig `json:"input_schema"`
}

type ToolInputSchemaConfig struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required"`
}

func validateTools(tools []ToolDefinition) error {
	if len(tools) == 0 {
		return fmt.Errorf("tools must define at least one tool")
	}

	names := make(map[string]struct{}, len(tools))
	for _, tool := range tools {
		if tool.Name == "" {
			return fmt.Errorf("tool name is required")
		}
		if tool.Description == "" {
			return fmt.Errorf("tool description is required for %q", tool.Name)
		}
		if tool.InputSchema.Type == "" {
			return fmt.Errorf("tool input_schema.type is required for %q", tool.Name)
		}
		if _, exists := names[tool.Name]; exists {
			return fmt.Errorf("duplicate tool name: %q", tool.Name)
		}
		names[tool.Name] = struct{}{}
	}

	return nil
}

func (t ToolDefinition) ToMCPTool() mcp.Tool {
	return mcp.Tool{
		Name:        t.Name,
		Description: t.Description,
		InputSchema: mcp.ToolInputSchema{
			Type:       t.InputSchema.Type,
			Properties: t.InputSchema.Properties,
			Required:   t.InputSchema.Required,
		},
	}
}
