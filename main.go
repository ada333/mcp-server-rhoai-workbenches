package main

import (
	"context"
	"log"

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

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Pods",
		Description: "list the pods in a namespace",
	}, tools.ListPods)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Workbenches",
		Description: "list the workbenches in a given project namespace",
	}, tools.ListWorkbenches)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List All Workbenches",
		Description: "list the workbenches across all namespaces",
	}, tools.ListAllWorkbenches)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Change Workbench Status",
		Description: "change the status of a workbench with given namein a given project namespace",
	}, tools.ChangeWorkbenchStatus)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "List Images",
		Description: "list the images in a given project namespace",
	}, tools.ListImages)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create Custom Image",
		Description: "create a new custom notebook image",
	}, tools.CreateCustomImage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create Workbench",
		Description: "create a new workbench with given name, image and image URL in a given project namespace",
	}, tools.CreateWorkbench)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Delete Workbench",
		Description: "delete a workbench with given name in a given project namespace",
	}, tools.DeleteWorkbench)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Delete Image",
		Description: "delete an image with given name",
	}, tools.DeleteImage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Create Hardware Profile",
		Description: "create a hardware profile with given name, description and resources",
	}, tools.CreateHardwareProfile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "Delete Hardware Profile",
		Description: "delete a hardware profile with given name",
	}, tools.DeleteHardwareProfile)

	server.AddResource(&mcp.Resource{
		URI:         "resource://mcp-server-rhoai/images",
		Name:        "Image Catalog",
		Description: "List of available notebook images their URLs and tags",
		MIMEType:    "application/json",
	}, resources.ImagesResourceHandler)

	server.AddPrompt(&mcp.Prompt{
		Name:        "create-workbench-prompt",
		Description: "Guide to create a workbench",
		Arguments: []*mcp.PromptArgument{
			{Name: "namespace", Description: "Target namespace", Required: true},
			{Name: "name", Description: "Workbench name", Required: true},
		},
	}, prompts.CreateWorkbenchPromptHandler)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
