package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func GetDefaultHardwareProfile() core.HardwareProfile {
	return core.HardwareProfile{
		HardwareProfileName: "default-profile",
		Resources: []core.HardwareProfileResource{
			{
				ResourceName:       "CPU",
				ResourceIdentifier: "cpu",
				ResourceType:       "CPU",
				DefaultCount:       "2",
				MinCount:           "1",
				MaxCount:           "4",
			},
			{
				ResourceName:       "Memory",
				ResourceIdentifier: "memory",
				ResourceType:       "Memory",
				DefaultCount:       "4Gi",
				MinCount:           "2Gi",
				MaxCount:           "8Gi",
			},
		},
	}
}

func DefaultHardwareResourceHandler(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	resources := GetDefaultHardwareProfile()

	jsonBytes, err := json.Marshal(resources)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default hardware resources: %v", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonBytes),
			},
		},
	}, nil
}
