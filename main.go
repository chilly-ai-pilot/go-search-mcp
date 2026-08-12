package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chilly-ai-pilot/go-search-mcp/config"
	"github.com/chilly-ai-pilot/go-search-mcp/handlers"
	_ "github.com/chilly-ai-pilot/go-search-mcp/handlers/web_search"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	configPath := flag.String("config", "", "path to MCP config file")
	flag.Parse()

	resolvedConfigPath, err := config.ResolveConfigPath(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve config path: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.Load(resolvedConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	s := server.NewMCPServer(
		cfg.Server.Name,
		cfg.Server.Version,
		server.WithToolCapabilities(true),
	)

	registry, err := handlers.NewRegistry(cfg, cfg.Tools)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize handlers: %v\n", err)
		os.Exit(1)
	}
	for _, toolDef := range cfg.Tools {
		handler, err := registry.Get(toolDef.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to register tool %q: %v\n", toolDef.Name, err)
			os.Exit(1)
		}
		s.AddTool(toolDef.ToMCPTool(), handler)
	}

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "stdio server error: %v\n", err)
		os.Exit(1)
	}
}
