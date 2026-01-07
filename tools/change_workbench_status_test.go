package tools

import (
	"context"
	"testing"
	"time"

	core "github.com/amaly/mcp-server-rhoai/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

func TestChangeWorkbenchStatus(t *testing.T) {
	orig := GetDynamicClient
	defer func() { GetDynamicClient = orig }()

	scheme := runtime.NewScheme()
	stoppedWorkbench := newUnstructuredWorkbench("StoppedWorkbench", "ns1")
	stoppedWorkbench.SetAnnotations(map[string]string{
		"kubeflow-resource-stopped": time.Now().UTC().Format(time.RFC3339),
	})
	runningWorkbench := newUnstructuredWorkbench("RunningWorkbench", "ns1")

	client := dynamicfake.NewSimpleDynamicClient(scheme,
		stoppedWorkbench,
		runningWorkbench,
	)

	GetDynamicClient = func() (dynamic.Interface, error) {
		return client, nil
	}

	_, out, err := ChangeWorkbenchStatus(context.Background(), nil, core.ChangeWorkbenchStatusInput{Namespace: "ns1", WorkbenchName: "StoppedWorkbench", Status: core.Stopped})
	if err != nil {
		t.Fatalf("ChangeWorkbenchStatus returned error: %v", err)
	}
	if out.Message != "Workbench StoppedWorkbench is already stopped" {
		t.Errorf("expected StoppedWorkbench is already stopped, got: %q", out.Message)
	}

	_, out, err = ChangeWorkbenchStatus(context.Background(), nil, core.ChangeWorkbenchStatusInput{Namespace: "ns1", WorkbenchName: "RunningWorkbench", Status: core.Running})
	if err != nil {
		t.Fatalf("ChangeWorkbenchStatus returned error: %v", err)
	}
	if out.Message != "Workbench RunningWorkbench is already running" {
		t.Errorf("expected RunningWorkbench is already running, got: %q", out.Message)
	}

	_, out, err = ChangeWorkbenchStatus(context.Background(), nil, core.ChangeWorkbenchStatusInput{Namespace: "ns1", WorkbenchName: "StoppedWorkbench", Status: core.Running})
	if err != nil {
		t.Fatalf("ChangeWorkbenchStatus returned error: %v", err)
	}

	if out.Message != "Workbench StoppedWorkbench is running" {
		t.Errorf("expected StoppedWorkbench is running, got: %q", out.Message)
	}

	_, out, err = ChangeWorkbenchStatus(context.Background(), nil, core.ChangeWorkbenchStatusInput{Namespace: "ns1", WorkbenchName: "RunningWorkbench", Status: core.Stopped})
	if err != nil {
		t.Fatalf("ChangeWorkbenchStatus returned error: %v", err)
	}
	if out.Message != "Workbench RunningWorkbench is stopped" {
		t.Errorf("expected RunningWorkbench is stopped, got: %q", out.Message)
	}
}
