package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chilly-ai-pilot/go-search-mcp/config"
)

func TestSearXNGProviderSearch(t *testing.T) {
	const totalResults = 5
	results := make([]map[string]interface{}, totalResults)
	for i := range results {
		results[i] = map[string]interface{}{
			"title":   "Result",
			"url":     "https://example.com",
			"content": "snippet",
			"score":   float64(i),
		}
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "test query" {
			t.Errorf("unexpected query: %q", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("format") != "json" {
			t.Errorf("unexpected format: %q", r.URL.Query().Get("format"))
		}
		if r.URL.Query().Get("categories") != "general" {
			t.Errorf("unexpected categories: %q", r.URL.Query().Get("categories"))
		}

		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"results": results,
		})
	}))
	defer server.Close()

	cfg := &config.WebSearchConfig{
		Provider:       "searxng",
		TimeoutSeconds: 5,
		APIURL:         server.URL,
		Providers: map[string]json.RawMessage{
			"searxng": json.RawMessage(`{"categories":"general"}`),
		},
	}

	provider, err := NewSearXNGProvider(cfg)
	if err != nil {
		t.Fatalf("NewSearXNGProvider: %v", err)
	}

	got, err := provider.Search(context.Background(), "test query", 2)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0].Title != "Result" || got[0].URL != "https://example.com" {
		t.Fatalf("unexpected first result: %+v", got[0])
	}
}

func TestTavilyProviderSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Errorf("unexpected auth header: %q", auth)
		}

		var req tavilyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Query != "test query" {
			t.Errorf("unexpected query: %q", req.Query)
		}
		if req.MaxResults != 2 {
			t.Errorf("unexpected max_results: %d", req.MaxResults)
		}
		if req.SearchDepth != "basic" {
			t.Errorf("unexpected search_depth: %q", req.SearchDepth)
		}

		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []map[string]interface{}{
				{"title": "Tavily Result", "url": "https://example.com", "content": "content", "score": 0.9},
			},
		})
	}))
	defer server.Close()

	t.Setenv("TAVILY_API_KEY", "test-key")

	cfg := &config.WebSearchConfig{
		Provider:       "tavily",
		TimeoutSeconds: 5,
		APIURL:         server.URL,
		APIKeyEnv:      "TAVILY_API_KEY",
		Providers: map[string]json.RawMessage{
			"tavily": json.RawMessage(`{"search_depth":"basic"}`),
		},
	}

	provider, err := NewTavilyProvider(cfg)
	if err != nil {
		t.Fatalf("NewTavilyProvider: %v", err)
	}

	got, err := provider.Search(context.Background(), "test query", 2)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Title != "Tavily Result" {
		t.Fatalf("unexpected result: %+v", got[0])
	}
}
