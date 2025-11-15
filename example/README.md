# Tea Shop MCP Server

A Model Context Protocol (MCP) server implementation in Go that provides a tea shop interface.

## Features

- Tea inventory management with search and stock checking
- 4 available tools: `list_teas`, `get_tea`, `search_teas`, `check_stock`
- Tea catalog resource: `teashop://catalog`

## Quick Start

```bash
# Build and run
go build -o bin/teashop-mcp-server .
./bin/teashop-mcp-server
```

## Claude Desktop Integration

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "teashop": {
      "command": "/path/to/bin/teashop-mcp-server",
      "args": []
    }
  }
}
```

## Usage Examples

- "List all available teas"
- "Search for green teas"
- "Get details about tea with ID 1"
- "Check if Royal Pu-erh is in stock"