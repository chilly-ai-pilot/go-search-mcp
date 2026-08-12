# Go Search MCP

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![MCP](https://img.shields.io/badge/MCP-Compatible-FF6B6B)](https://modelcontextprotocol.io)

一个基于 MCP (Model Context Protocol) 的搜索服务，支持多种搜索引擎后端，为 AI 助手提供联网搜索能力。

## ✨ 特性

- 🔍 **多搜索引擎支持**：Tavily、SearXNG，可轻松扩展
- 🧩 **MCP 协议兼容**：可作为 MCP Server 与 Claude Desktop 等客户端集成
- ⚙️ **配置驱动**：通过 JSON 配置文件管理所有设置
- 🔌 **插件化架构**：新增搜索引擎无需修改核心代码
- 🚀 **轻量高效**：纯 Go 实现，单二进制文件部署

## 📦 支持的搜索引擎

| Provider | 类型 | 说明 |
|----------|------|------|
| [Tavily](https://tavily.com) | 商业 API | 需要 API Key |
| [SearXNG](https://docs.searxng.org) | 自托管 | 需要部署 SearXNG 实例 |

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Tavily API Key（使用 Tavily 时）或 SearXNG 实例（使用 SearXNG 时）

### 安装

```bash
# 克隆项目
git clone https://github.com/chilly-ai-pilot/go-search-mcp.git
cd go-search-mcp

# 下载依赖
go mod tidy

# 构建
go build -ldflags="-s -w" -o go-search-mcp
```

### 配置

#### 使用 Tavily

创建配置文件 `mcp-config.json`：

```json
{
  "server": {
    "name": "go-search-mcp",
    "version": "1.0.0"
  },
  "handlers": {
    "web_search": {
      "provider": "tavily",
      "max_results": 5,
      "timeout_seconds": 30,
      "api_url": "https://api.tavily.com/search",
      "api_key_env": "TAVILY_API_KEY",
      "providers": {
        "tavily": {
          "search_depth": "basic"
        }
      }
    }
  },
  "tools": [
    {
      "name": "web_search",
      "description": "当知识库内容不足以回答用户问题时，上网搜索补充信息。",
      "input_schema": {
        "type": "object",
        "properties": {
          "query": {
            "type": "string",
            "description": "搜索关键词"
          },
          "max_results": {
            "type": "integer",
            "description": "返回结果的最大数量"
          }
        },
        "required": ["query"]
      }
    }
  ]
}
```

#### 使用 SearXNG

```json
{
  "handlers": {
    "web_search": {
      "provider": "searxng",
      "max_results": 5,
      "timeout_seconds": 30,
      "api_url": "http://localhost:8080/search",
      "providers": {
        "searxng": {
          "categories": "general"
        }
      }
    }
  }
}
```

### 运行

```bash
# 设置环境变量（使用 Tavily 时需要）
export TAVILY_API_KEY="your-api-key"

# 启动服务
./go-search-mcp -config mcp-config.json
```

## 🔌 与 Claude Desktop 集成

在 Claude Desktop 的配置文件（`~/Library/Application Support/Claude/claude_desktop_config.json`）中添加：

```json
{
  "mcpServers": {
    "go-search-mcp": {
      "command": "/path/to/go-search-mcp",
      "args": [
        "-config",
        "/path/to/mcp-config.json"
      ],
      "env": {
        "TAVILY_API_KEY": "your-api-key"
      }
    }
  }
}
```

## 🏗️ 架构设计

```
go-search-mcp/
├── main.go                 # 入口文件
├── config/                 # 配置管理
│   ├── config.go          # 配置加载与验证
│   └── tool_config.go     # Tool 定义
├── handlers/               # 处理器（可扩展）
│   ├── registry.go        # Handler 注册中心
│   └── web_search/        # Web 搜索处理器
│       ├── handler.go     # 处理逻辑
│       ├── providers/     # 搜索引擎实现
│       │   ├── provider.go    # Provider 接口与注册
│       │   ├── tavily.go      # Tavily 实现
│       │   └── searxng.go     # SearXNG 实现
│       └── models/        # 数据模型
└── mcp-config.json        # 配置文件示例
```

### 开闭原则设计

- **新增搜索引擎**：在 `providers/` 下新建文件，实现 `SearchProvider` 接口，通过 `init()` 注册
- **新增 Tool**：在 `handlers/` 下新建目录，实现处理逻辑，通过 `init()` 注册
- **无需修改核心代码**：所有扩展通过注册机制完成

## 🧪 测试

```bash
go test ./...
```

## 📝 配置说明

### 通用配置

| 字段 | 类型 | 说明 |
|------|------|------|
| `server.name` | string | 服务名称 |
| `server.version` | string | 服务版本 |
| `handlers.web_search.provider` | string | 搜索引擎：`tavily` 或 `searxng` |
| `handlers.web_search.max_results` | int | 默认返回结果数 |
| `handlers.web_search.timeout_seconds` | int | 请求超时时间 |
| `handlers.web_search.api_url` | string | API 地址 |
| `handlers.web_search.api_key_env` | string | API Key 环境变量名（仅 Tavily） |
| `handlers.web_search.providers` | object | 各搜索引擎专属配置 |

### Provider 专属配置

#### Tavily
```json
{
  "tavily": {
    "search_depth": "basic"  // 或 "advanced"
  }
}
```

#### SearXNG
```json
{
  "searxng": {
    "categories": "general"  // 搜索分类
  }
}
```

## 🛠️ 扩展开发

### 添加新的搜索引擎

1. 在 `handlers/web_search/providers/` 下创建 `xxx.go`：

```go
package providers

import (
    "context"
    "github.com/chilly-ai-pilot/go-search-mcp/config"
    "github.com/chilly-ai-pilot/go-search-mcp/handlers/web_search/models"
)

const xxxProviderType = "xxx"

func init() {
    RegisterProvider(xxxProviderType, func(cfg *config.WebSearchConfig) (SearchProvider, error) {
        return NewXXXProvider(cfg)
    })
}

type XXXProvider struct { /* ... */ }

func (p *XXXProvider) Search(ctx context.Context, query string, maxResults int) ([]models.SearchResult, error) {
    // 实现搜索逻辑
}
```

2. 在配置文件的 `providers` 中添加专属配置：

```json
{
  "providers": {
    "xxx": {
      "custom_field": "value"
    }
  }
}
```

### 添加新的 Tool

1. 在 `handlers/` 下新建目录，实现 `HandlerBuilder` 接口
2. 在 `tools` 数组中添加 Tool 定义
3. 在 `main.go` 中导入新 handler

## 📄 License

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 🙏 致谢

- [mcp-go](https://github.com/mark3labs/mcp-go) - MCP Go SDK
- [Tavily](https://tavily.com) - 搜索 API
- [SearXNG](https://searxng.org) - 开源元搜索引擎
