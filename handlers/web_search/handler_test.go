package websearch

import (
	"context"
	"testing"

	"github.com/chilly-ai-pilot/go-search-mcp/handlers/web_search/models"
	"github.com/mark3labs/mcp-go/mcp"
)

type fakeSearchProvider struct {
	gotMaxResults int
}

func (f *fakeSearchProvider) Search(ctx context.Context, query string, maxResults int) ([]models.SearchResult, error) {
	f.gotMaxResults = maxResults
	return []models.SearchResult{{Title: "result", URL: "https://example.com", Content: "snippet", Score: 1.0}}, nil
}

func TestHandleClampsMaxResultsToConfiguredMax(t *testing.T) {
	provider := &fakeSearchProvider{}
	h := NewHandler(provider, 2, 5)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"query":       "test query",
				"max_results": float64(10),
			},
		},
	}

	result, err := h.Handle(context.Background(), request)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if provider.gotMaxResults != 5 {
		t.Fatalf("expected max_results to be clamped to 5, got %d", provider.gotMaxResults)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
