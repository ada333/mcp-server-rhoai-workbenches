package main

import (
	"context"
	"log"
	"os"

	"github.com/amaly/mcp-server-rhoai/prompts"
	"github.com/amaly/mcp-server-rhoai/resources"
	"github.com/amaly/mcp-server-rhoai/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "workbencheslist",
		Version: "v1.0.0",
	}, nil)

	mode := os.Getenv("MCP_RHOAI_MODE")
	if mode == "write" {
		tools.RegisterWriteTools(server)
	}
	tools.RegisterReadOnlyTools(server)

	resources.RegisterAllResources(server)
	prompts.RegisterAllPrompts(server)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
