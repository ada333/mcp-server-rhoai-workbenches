package tools

import (
	"context"
	"fmt"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteImage(ctx context.Context, req *mcp.CallToolRequest, input core.DeleteImageInput) (*mcp.CallToolResult, core.WorkbenchOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.WorkbenchOutput{}, err
	}

	err = dyn.Resource(core.ImagesGVR).Namespace(core.GetDefaultNamespace()).Delete(ctx, input.ImageName, metav1.DeleteOptions{})
	if err != nil {
		return nil, core.WorkbenchOutput{}, fmt.Errorf("failed to delete image %s: %v", input.ImageName, err)
	}

	return nil, core.WorkbenchOutput{Message: fmt.Sprintf("Image %s was successfully deleted", input.ImageName)}, nil
}
