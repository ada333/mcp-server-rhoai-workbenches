package prompts

import "github.com/modelcontextprotocol/go-sdk/mcp"

func RegisterAllPrompts(server *mcp.Server) {
	server.AddPrompt(&mcp.Prompt{
		Name:        "create-workbench-prompt",
		Description: "Guide to create a workbench",
		Arguments: []*mcp.PromptArgument{
			{Name: "namespace", Description: "Target namespace", Required: true},
			{Name: "name", Description: "Workbench name", Required: true},
		},
	}, CreateWorkbenchPromptHandler)
}
