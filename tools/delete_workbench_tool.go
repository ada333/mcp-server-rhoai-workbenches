package tools

import (
	"context"
	"fmt"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteWorkbench(ctx context.Context, req *mcp.CallToolRequest, input core.DeleteWorkbenchInput) (*mcp.CallToolResult, core.WorkbenchOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.WorkbenchOutput{}, err
	}

	err = dyn.Resource(core.WorkbenchesGVR).Namespace(input.Namespace).Delete(ctx, input.WorkbenchName, metav1.DeleteOptions{})
	if err != nil {
		return nil, core.WorkbenchOutput{}, fmt.Errorf("failed to delete workbench %s: %v", input.WorkbenchName, err)
	}

	return nil, core.WorkbenchOutput{Message: fmt.Sprintf("Workbench %s was successfully deleted", input.WorkbenchName)}, nil
}
