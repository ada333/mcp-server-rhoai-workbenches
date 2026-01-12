package tools

import (
	"context"
	"fmt"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteHardwareProfile(ctx context.Context, req *mcp.CallToolRequest, input core.DeleteHardwareProfileInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	err = dyn.Resource(core.HardwareProfilesGVR).Namespace(core.GetDefaultNamespace()).Delete(ctx, input.HardwareProfileName, metav1.DeleteOptions{})
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to delete hardware profile %s: %v", input.HardwareProfileName, err)
	}

	return nil, core.DefaultToolOutput{Message: fmt.Sprintf("Hardware Profile %s was successfully deleted", input.HardwareProfileName)}, nil
}
