package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chilly-ai-pilot/go-search-mcp/config"
	"github.com/chilly-ai-pilot/go-search-mcp/handlers/web_search/models"
)

type SearchProvider interface {
	Search(ctx context.Context, query string, maxResults int) ([]models.SearchResult, error)
}

type factory func(*config.WebSearchConfig) (SearchProvider, error)

var registry = map[string]factory{}

func RegisterProvider(name string, create factory) {
	registry[name] = create
}

func New(cfg *config.WebSearchConfig) (SearchProvider, error) {
	create, ok := registry[cfg.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported search provider: %s", cfg.Provider)
	}
	return create(cfg)
}

type searchResultItem struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

type searchResultsResponse struct {
	Results []searchResultItem `json:"results"`
}

func newHTTPClient(cfg *config.WebSearchConfig) *http.Client {
	return &http.Client{Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second}
}

func doRequest(ctx context.Context, client *http.Client, req *http.Request, providerName string) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", providerName, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s response: %w", providerName, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned status %d: %s", providerName, resp.StatusCode, string(body))
	}
	return body, nil
}

func parseSearchResults(body []byte, providerName string) ([]models.SearchResult, error) {
	var resp searchResultsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse %s response: %w", providerName, err)
	}
	return toSearchResults(resp.Results), nil
}

func toSearchResults(items []searchResultItem) []models.SearchResult {
	results := make([]models.SearchResult, 0, len(items))
	for _, item := range items {
		results = append(results, models.SearchResult{
			Title:   item.Title,
			URL:     item.URL,
			Content: item.Content,
			Score:   item.Score,
		})
	}
	return results
}

func limitResults(results []models.SearchResult, maxResults int) []models.SearchResult {
	if maxResults > 0 && len(results) > maxResults {
		return results[:maxResults]
	}
	return results
}
