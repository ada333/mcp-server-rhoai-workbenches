package tools

import (
	"context"
	"fmt"

	"github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func CreatePVCTool(ctx context.Context, req *mcp.CallToolRequest, input core.CreatePVCInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	err = createPersistentVolumeClaim(ctx, dyn, input.Namespace, input.PVCName, input.Size)
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to create PVC: %v", err)
	}

	return nil, core.DefaultToolOutput{Message: "PVC was succesfully created!"}, nil
}
