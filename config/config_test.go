package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRelativeConfigPath(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "mcp-config.json")

	configJSON := `{
  "server": {"name": "test", "version": "1.0.0"},
  "handlers": {
    "web_search": {
      "provider": "searxng",
      "max_results": 5,
      "timeout_seconds": 30,
      "api_url": "http://localhost:8888/search",
      "providers": {
        "tavily": {"search_depth": "basic"},
        "searxng": {"categories": "general"}
      }
    }
  },
  "tools": [{
    "name": "web_search",
    "description": "test",
    "input_schema": {
      "type": "object",
      "properties": {"query": {"type": "string"}},
      "required": ["query"]
    }
  }]
}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cfg, err := Load("mcp-config.json")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Tools) != 1 || cfg.Tools[0].Name != "web_search" {
		t.Fatalf("unexpected tools: %+v", cfg.Tools)
	}
}

func TestLoadWebSearchDefaultResultsNormalizesToMax(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "mcp-config.json")

	configJSON := `{
  "server": {"name": "test", "version": "1.0.0"},
  "handlers": {
    "web_search": {
      "provider": "searxng",
      "max_results": 5,
      "default_results": 10,
      "timeout_seconds": 30,
      "api_url": "http://localhost:8888/search",
      "providers": {
        "tavily": {"search_depth": "basic"},
        "searxng": {"categories": "general"}
      }
    }
  },
  "tools": [{
    "name": "web_search",
    "description": "test",
    "input_schema": {
      "type": "object",
      "properties": {"query": {"type": "string"}},
      "required": ["query"]
    }
  }]
}`
	if err := os.WriteFile(configPath, []byte(configJSON), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Handlers.WebSearch.DefaultResults != cfg.Handlers.WebSearch.MaxResults {
		t.Fatalf("expected default_results normalized to max_results, got %d", cfg.Handlers.WebSearch.DefaultResults)
	}
}
