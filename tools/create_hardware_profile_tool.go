package tools

import (
	"context"
	"fmt"
	"time"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CreateHardwareProfile(ctx context.Context, req *mcp.CallToolRequest, input core.CreateHardwareProfileInput) (*mcp.CallToolResult, core.DefaultToolOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.DefaultToolOutput{}, err
	}

	namespace := core.GetDefaultNamespace()

	identifiers := make([]interface{}, len(input.Resources))
	for i, res := range input.Resources {
		identifiers[i] = map[string]interface{}{
			"displayName":  res.ResourceName,
			"identifier":   res.ResourceIdentifier,
			"resourceType": res.ResourceType,
			"defaultCount": res.DefaultCount,
			"maxCount":     res.MaxCount,
			"minCount":     res.MinCount,
		}
	}

	annotations := map[string]interface{}{
		"opendatahub.io/dashboard-feature-visibility": "[]",
		"opendatahub.io/disabled":                     "false",
		"opendatahub.io/display-name":                 input.HardwareProfileName,
		"opendatahub.io/modified-date":                time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}

	hardwareProfile := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "infrastructure.opendatahub.io/v1",
			"kind":       "HardwareProfile",
			"metadata": map[string]interface{}{
				"name":        input.HardwareProfileName,
				"namespace":   namespace,
				"annotations": annotations,
			},
			"spec": map[string]interface{}{
				"identifiers": identifiers,
			},
		},
	}

	_, err = dyn.Resource(core.HardwareProfilesGVR).Namespace(namespace).Create(ctx, hardwareProfile, metav1.CreateOptions{})
	if err != nil {
		return nil, core.DefaultToolOutput{}, fmt.Errorf("failed to create hardware profile: %v", err)
	}

	return nil, core.DefaultToolOutput{
		Message: "Hardware Profile was successfully created!",
	}, nil
}
