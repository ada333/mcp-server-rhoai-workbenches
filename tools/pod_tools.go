package tools

import (
	"context"
	"fmt"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListPods(ctx context.Context, req *mcp.CallToolRequest, input core.ListWorkbenchesInput) (*mcp.CallToolResult, core.PodsOutput, error) {
	clientset, err := GetClientSet()
	if err != nil {
		return nil, core.PodsOutput{}, err
	}

	pods, err := clientset.CoreV1().Pods(input.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, core.PodsOutput{}, fmt.Errorf("failed to list pods: %v", err)
	}

	msg := ""
	for _, pod := range pods.Items {
		msg += fmt.Sprintf("- %s (%s)\n", pod.Name, pod.Status.Phase)
	}
	return nil, core.PodsOutput{Pods: msg}, nil
}
