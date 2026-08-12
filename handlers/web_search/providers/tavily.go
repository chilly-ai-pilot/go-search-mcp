package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/chilly-ai-pilot/go-search-mcp/config"
	"github.com/chilly-ai-pilot/go-search-mcp/handlers/web_search/models"
)

const tavilyProviderType = "tavily"

func init() {
	RegisterProvider(tavilyProviderType, func(cfg *config.WebSearchConfig) (SearchProvider, error) {
		return NewTavilyProvider(cfg)
	})
}

type tavilyConfig struct {
	SearchDepth string `json:"search_depth"`
}

type TavilyProvider struct {
	apiURL      string
	apiKey      string
	searchDepth string
	client      *http.Client
}

func NewTavilyProvider(cfg *config.WebSearchConfig) (*TavilyProvider, error) {
	var providerCfg tavilyConfig
	if err := cfg.DecodeActiveProvider(&providerCfg); err != nil {
		return nil, err
	}

	apiKey := os.Getenv(cfg.APIKeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("environment variable %s is not set", cfg.APIKeyEnv)
	}

	return &TavilyProvider{
		apiURL:      cfg.APIURL,
		apiKey:      apiKey,
		searchDepth: providerCfg.SearchDepth,
		client:      newHTTPClient(cfg),
	}, nil
}

type tavilyRequest struct {
	Query       string `json:"query"`
	MaxResults  int    `json:"max_results"`
	SearchDepth string `json:"search_depth,omitempty"`
}

func (p *TavilyProvider) Search(ctx context.Context, query string, maxResults int) ([]models.SearchResult, error) {
	body, err := json.Marshal(tavilyRequest{
		Query:       query,
		MaxResults:  maxResults,
		SearchDepth: p.searchDepth,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal tavily request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create tavily request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	respBody, err := doRequest(ctx, p.client, req, tavilyProviderType)
	if err != nil {
		return nil, err
	}

	return parseSearchResults(respBody, tavilyProviderType)
}
