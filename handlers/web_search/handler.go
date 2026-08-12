package websearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chilly-ai-pilot/go-search-mcp/config"
	"github.com/chilly-ai-pilot/go-search-mcp/handlers"
	"github.com/chilly-ai-pilot/go-search-mcp/handlers/web_search/providers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const HandlerType = "web_search"

func init() {
	handlers.RegisterHandler(HandlerType, func(cfg *config.Config) (server.ToolHandlerFunc, error) {
		provider, err := providers.New(&cfg.Handlers.WebSearch)
		if err != nil {
			return nil, err
		}
		return NewHandler(provider, cfg.Handlers.WebSearch.MaxResults).Handle, nil
	})
}

type Handler struct {
	provider   providers.SearchProvider
	maxResults int
}

func NewHandler(provider providers.SearchProvider, maxResults int) *Handler {
	return &Handler{
		provider:   provider,
		maxResults: maxResults,
	}
}

func (h *Handler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("参数类型错误"), nil
	}

	query, ok := args["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("缺少必填参数 query"), nil
	}

	maxResults := h.maxResults
	if raw, ok := args["max_results"]; ok {
		switch v := raw.(type) {
		case float64:
			maxResults = int(v)
		case int:
			maxResults = v
		}
	}

	results, err := h.provider.Search(ctx, query, maxResults)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("搜索失败: %v", err)), nil
	}

	bytes, err := json.Marshal(results)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("结果序列化失败: %v", err)), nil
	}

	return mcp.NewToolResultText(string(bytes)), nil
}
