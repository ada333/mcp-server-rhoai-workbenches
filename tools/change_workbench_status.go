package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	core "github.com/ada333/MCP-test/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

func IsWorkbenchStopped(ctx context.Context, dyn dynamic.Interface, namespace, workbenchName string) (bool, error) {
	current, err := dyn.Resource(core.WorkbenchesGVR).Namespace(namespace).Get(ctx, workbenchName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get workbench %s: %v", workbenchName, err)
	}
	currentAnnotations := current.GetAnnotations()
	currentStopped := false
	if currentAnnotations != nil {
		if _, ok := currentAnnotations["kubeflow-resource-stopped"]; ok {
			currentStopped = true
		}
	}
	return currentStopped, nil
}

func ChangeWorkbenchStatus(ctx context.Context, req *mcp.CallToolRequest, input core.ChangeWorkbenchStatusInput) (*mcp.CallToolResult, core.WorkbenchOutput, error) {
	dyn, err := GetDynamicClient()
	if err != nil {
		return nil, core.WorkbenchOutput{}, err
	}

	stopped, err := IsWorkbenchStopped(ctx, dyn, input.Namespace, input.WorkbenchName)
	if err != nil {
		return nil, core.WorkbenchOutput{}, err
	}
	if (input.Status == core.Stopped && stopped) || (input.Status == core.Running && !stopped) {
		return nil, core.WorkbenchOutput{Message: fmt.Sprintf("Workbench %s is already %s", input.WorkbenchName, input.Status)}, nil
	}

	patch := map[string]interface{}{}
	annotations := map[string]interface{}{}
	if input.Status == core.Stopped {
		annotations["kubeflow-resource-stopped"] = time.Now().UTC().Format(time.RFC3339)
	} else {
		annotations["kubeflow-resource-stopped"] = nil
	}
	patch["metadata"] = map[string]interface{}{
		"annotations": annotations,
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return nil, core.WorkbenchOutput{}, fmt.Errorf("failed to marshal patch: %v", err)
	}

	_, err = dyn.Resource(core.WorkbenchesGVR).Namespace(input.Namespace).Patch(
		ctx,
		input.WorkbenchName,
		k8stypes.MergePatchType,
		patchBytes,
		metav1.PatchOptions{},
	)
	if err != nil {
		return nil, core.WorkbenchOutput{}, fmt.Errorf("failed to %s workbench %s: %v", input.Status, input.WorkbenchName, err)
	}

	return nil, core.WorkbenchOutput{Message: fmt.Sprintf("Workbench %s is %s", input.WorkbenchName, input.Status)}, nil
}
