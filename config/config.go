package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Server   ServerConfig     `json:"server"`
	Handlers HandlersConfig   `json:"handlers"`
	Tools    []ToolDefinition `json:"tools"`
}

type HandlersConfig struct {
	WebSearch WebSearchConfig `json:"web_search"`
}

type ServerConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type WebSearchConfig struct {
	Provider       string                     `json:"provider"`
	MaxResults     int                        `json:"max_results"`
	DefaultResults int                        `json:"default_results"`
	TimeoutSeconds int                        `json:"timeout_seconds"`
	APIURL         string                     `json:"api_url"`
	APIKeyEnv      string                     `json:"api_key_env"`
	Providers      map[string]json.RawMessage `json:"providers"`
}

func (s *WebSearchConfig) DecodeProvider(name string, v any) error {
	raw, ok := s.Providers[name]
	if !ok || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, v); err != nil {
		return fmt.Errorf("decode handlers.web_search.providers.%s: %w", name, err)
	}
	return nil
}

func (s *WebSearchConfig) DecodeActiveProvider(v any) error {
	return s.DecodeProvider(s.Provider, v)
}

func Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config path is required")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve config path %s: %w", path, err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", absPath, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", absPath, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Server.Name == "" {
		return fmt.Errorf("server.name is required")
	}
	if c.Server.Version == "" {
		return fmt.Errorf("server.version is required")
	}
	if c.Handlers.WebSearch.Provider == "" {
		return fmt.Errorf("handlers.web_search.provider is required")
	}
	if len(c.Handlers.WebSearch.Providers) == 0 {
		return fmt.Errorf("handlers.web_search.providers is required")
	}
	if _, ok := c.Handlers.WebSearch.Providers[c.Handlers.WebSearch.Provider]; !ok {
		return fmt.Errorf("handlers.web_search.providers.%s is required", c.Handlers.WebSearch.Provider)
	}
	if c.Handlers.WebSearch.MaxResults <= 0 {
		return fmt.Errorf("handlers.web_search.max_results must be greater than 0")
	}
	if c.Handlers.WebSearch.DefaultResults <= 0 {
		c.Handlers.WebSearch.DefaultResults = c.Handlers.WebSearch.MaxResults
	}
	if c.Handlers.WebSearch.DefaultResults > c.Handlers.WebSearch.MaxResults {
		c.Handlers.WebSearch.DefaultResults = c.Handlers.WebSearch.MaxResults
	}
	if c.Handlers.WebSearch.TimeoutSeconds <= 0 {
		return fmt.Errorf("handlers.web_search.timeout_seconds must be greater than 0")
	}
	if c.Handlers.WebSearch.APIURL == "" {
		return fmt.Errorf("handlers.web_search.api_url is required")
	}

	return validateTools(c.Tools)
}

func ResolveConfigPath(flagPath string) (string, error) {
	path := flagPath
	if path == "" {
		path = os.Getenv("MCP_CONFIG_PATH")
	}
	if path == "" {
		return "", fmt.Errorf("config path is required: use -config or MCP_CONFIG_PATH")
	}
	return filepath.Abs(path)
}
