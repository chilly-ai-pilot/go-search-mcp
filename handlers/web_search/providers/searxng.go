package providers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/chilly-ai-pilot/go-search-mcp/config"
	"github.com/chilly-ai-pilot/go-search-mcp/handlers/web_search/models"
)

const searxngProviderType = "searxng"

func init() {
	RegisterProvider(searxngProviderType, func(cfg *config.WebSearchConfig) (SearchProvider, error) {
		return NewSearXNGProvider(cfg)
	})
}

type searxngConfig struct {
	Categories string `json:"categories"`
}

type SearXNGProvider struct {
	apiURL     string
	categories string
	client     *http.Client
}

func NewSearXNGProvider(cfg *config.WebSearchConfig) (*SearXNGProvider, error) {
	var providerCfg searxngConfig
	if err := cfg.DecodeActiveProvider(&providerCfg); err != nil {
		return nil, err
	}

	return &SearXNGProvider{
		apiURL:     cfg.APIURL,
		categories: providerCfg.Categories,
		client:     newHTTPClient(cfg),
	}, nil
}

func (p *SearXNGProvider) Search(ctx context.Context, query string, maxResults int) ([]models.SearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	if p.categories != "" {
		params.Set("categories", p.categories)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("create searxng request: %w", err)
	}

	respBody, err := doRequest(ctx, p.client, req, searxngProviderType)
	if err != nil {
		return nil, err
	}

	results, err := parseSearchResults(respBody, searxngProviderType)
	if err != nil {
		return nil, err
	}
	return limitResults(results, maxResults), nil
}
