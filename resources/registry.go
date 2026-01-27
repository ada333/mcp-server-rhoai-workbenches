package resources

import "github.com/modelcontextprotocol/go-sdk/mcp"

func RegisterAllResources(server *mcp.Server) {
	server.AddResource(&mcp.Resource{
		URI:         "resource://mcp-server-rhoai/images",
		Name:        "Image Catalog",
		Description: "List of available notebook images their URLs and tags",
		MIMEType:    "application/json",
	}, ImagesResourceHandler)

	server.AddResource(&mcp.Resource{
		URI:         "resource://mcp-server-rhoai/hardware-resources",
		Name:        "Hardware Resources",
		Description: "List of available hardware resources",
		MIMEType:    "application/json",
	}, DefaultHardwareResourceHandler)
}
