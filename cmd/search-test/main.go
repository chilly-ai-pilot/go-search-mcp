package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chilly-ai-pilot/go-search-mcp/config"
	"github.com/chilly-ai-pilot/go-search-mcp/handlers/web_search/providers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <mcp-config.json> [max_results]\n", os.Args[0])
		os.Exit(1)
	}

	cfgPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve config path: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	maxResults := cfg.Handlers.WebSearch.MaxResults
	if len(os.Args) >= 3 {
		if _, err := fmt.Sscanf(os.Args[2], "%d", &maxResults); err != nil {
			fmt.Fprintf(os.Stderr, "parse max_results: %v\n", err)
			os.Exit(1)
		}
	}

	provider, err := providers.New(&cfg.Handlers.WebSearch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init provider: %v\n", err)
		os.Exit(1)
	}

	results, err := provider.Search(context.Background(), "MCP protocol introduction", maxResults)
	if err != nil {
		fmt.Fprintf(os.Stderr, "search failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("provider=%s max_results=%d returned=%d\n", cfg.Handlers.WebSearch.Provider, maxResults, len(results))
	for i, r := range results {
		fmt.Printf("[%d] title=%q url=%q score=%.4f content_len=%d\n", i+1, r.Title, r.URL, r.Score, len(r.Content))
	}

	out, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(out))
}
