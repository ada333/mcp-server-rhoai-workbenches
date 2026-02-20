# mcp-server-RHOAI
This project aims to implement MCP servers working with RHOAI and is used as a bachelors thesis.

## Tool Modes

The server supports two modes:

- **Read-only mode** (default): Only listing and querying tools are available. This mode is safer and prevents accidental modifications to your cluster.
- **Read-write mode**: All tools are available, including create, delete, and modify operations.

### Configuring the Mode

The mode is controlled by the `MCP_RHOAI_MODE` environment variable:

- If not set or set to any value other than `write`: Read-only mode (default)
- If set to `write`: Read-write mode with full access to all tools

Example configuration for read-write mode in Cursor:
```json
{
  "mcpServers": {
    "mcp-server-rhoai": {
      "command": "/home/amaly/mcp-server-rhoai/mcp-server-rhoai",
      "env": {
        "MCP_RHOAI_MODE": "write"
      }
    }
  }
}
```

# Setting up MCP server in Cursor
- Compile the main.go function (in repo -> go build -o any_name)
- In Cursor press Ctrl + Shift + P and type in Open MCP Settings
- Click on New MCP server

Example code you can add to make new MCP server in read-only mode (the command is path to the binary you build):
```
{
  "mcpServers": {
    "mcp-server-rhoai": {
      "command": "/home/amaly/mcp-server-rhoai/mcp-server-rhoai"
    }
  }
}
```
- Check that server is enabled in Cursor MCP settings and has some tools you can use
- now you can use the mcp tools just by talking to AI agent in Cursor (example prompt: find me workbenches in namespace mcp-test)


(to use the tools operating with OpenShift cluster you need to be logged in, the permissions are according to your login)


## Linting

This repository uses golangci-lint.

- Install:
  - Using Go: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
  - Or see `https://golangci-lint.run` for other options.
- Run locally:
  - With Make: `make lint`
  - Directly: `golangci-lint run`

Configuration is in `.golangci.yml`.
