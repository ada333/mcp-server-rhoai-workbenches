// here should be unit tests for the tools
// in the tests we use mocking to inject fake clientset and dynamic client

package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	core "github.com/ada333/MCP-test/core"
	"github.com/ada333/MCP-test/resources"
	"github.com/ada333/MCP-test/tools"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestListPods_Success(t *testing.T) {
	orig := tools.GetClientSet
	defer func() { tools.GetClientSet = orig }()

	ns := "test-ns"
	client := fake.NewSimpleClientset(
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-a",
				Namespace: ns,
			},
			Status: corev1.PodStatus{Phase: corev1.PodRunning},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod-other",
				Namespace: "other-ns",
			},
			Status: corev1.PodStatus{Phase: corev1.PodSucceeded},
		},
	)
	tools.GetClientSet = func() (kubernetes.Interface, error) {
		return client, nil
	}

	_, out, err := tools.ListPods(context.Background(), nil, core.ListWorkbenchesInput{Namespace: ns})
	if err != nil {
		t.Fatalf("ListPods returned error: %v", err)
	}
	if !strings.Contains(out.Pods, "- pod-a (Running)\n") {
		t.Errorf("expected pod-a Running in output, got: %q", out.Pods)
	}
	if strings.Contains(out.Pods, "pod-other") {
		t.Errorf("did not expect pod-other in output, got: %q", out.Pods)
	}
}

func newUnstructuredWorkbench(name, namespace string) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(core.WorkbenchesGVR.GroupVersion().WithKind("Notebook"))
	u.SetName(name)
	u.SetNamespace(namespace)
	return u
}

func TestListWorkbenches(t *testing.T) {
	orig := tools.GetDynamicClient
	defer func() { tools.GetDynamicClient = orig }()

	ns := "test-ns"
	scheme := runtime.NewScheme()
	client := dynamicfake.NewSimpleDynamicClient(scheme,
		newUnstructuredWorkbench("wb-1", ns),
		newUnstructuredWorkbench("wb-other", "other-ns"),
	)

	tools.GetDynamicClient = func() (dynamic.Interface, error) {
		return client, nil
	}

	_, out, err := tools.ListWorkbenches(context.Background(), nil, core.ListWorkbenchesInput{Namespace: ns})
	if err != nil {
		t.Fatalf("ListWorkbenches returned error: %v", err)
	}

	if !strings.Contains(out.Workbenches, "- wb-1\n") {
		t.Errorf("expected wb-1 in output, got: %q", out.Workbenches)
	}
	if strings.Contains(out.Workbenches, "wb-other") {
		t.Errorf("did not expect wb-other in output, got: %q", out.Workbenches)
	}
}

func TestListAllWorkbenches(t *testing.T) {
	orig := tools.GetDynamicClient
	defer func() { tools.GetDynamicClient = orig }()

	scheme := runtime.NewScheme()
	client := dynamicfake.NewSimpleDynamicClient(scheme,
		newUnstructuredWorkbench("wb-1", "ns1"),
		newUnstructuredWorkbench("wb-2", "ns2"),
	)

	tools.GetDynamicClient = func() (dynamic.Interface, error) {
		return client, nil
	}

	_, out, err := tools.ListAllWorkbenches(context.Background(), nil, core.ListWorkbenchesInput{Namespace: ""})
	if err != nil {
		t.Fatalf("ListAllWorkbenches returned error: %v", err)
	}

	if !strings.Contains(out.Workbenches, "- wb-1\n") {
		t.Errorf("expected wb-1 in output, got: %q", out.Workbenches)
	}
	if !strings.Contains(out.Workbenches, "- wb-2\n") {
		t.Errorf("expected wb-2 in output, got: %q", out.Workbenches)
	}
}

// TODO
func TestChangeWorkbenchStatus(t *testing.T) {
	orig := tools.GetDynamicClient
	defer func() { tools.GetDynamicClient = orig }()

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

	tools.GetDynamicClient = func() (dynamic.Interface, error) {
		return client, nil
	}

	// tests all the combinations of status changes
	// maybe checking the annotations would be better than output message
	_, out, err := tools.ChangeWorkbenchStatus(context.Background(), nil, core.ChangeWorkbenchStatusInput{Namespace: "ns1", WorkbenchName: "StoppedWorkbench", Status: core.Stopped})
	if err != nil {
		t.Fatalf("ChangeWorkbenchStatus returned error: %v", err)
	}
	if out.Message != "Workbench StoppedWorkbench is already stopped" {
		t.Errorf("expected StoppedWorkbench is already stopped, got: %q", out.Message)
	}

	_, out, err = tools.ChangeWorkbenchStatus(context.Background(), nil, core.ChangeWorkbenchStatusInput{Namespace: "ns1", WorkbenchName: "RunningWorkbench", Status: core.Running})
	if err != nil {
		t.Fatalf("ChangeWorkbenchStatus returned error: %v", err)
	}
	if out.Message != "Workbench RunningWorkbench is already running" {
		t.Errorf("expected RunningWorkbench is already running, got: %q", out.Message)
	}

	_, out, err = tools.ChangeWorkbenchStatus(context.Background(), nil, core.ChangeWorkbenchStatusInput{Namespace: "ns1", WorkbenchName: "StoppedWorkbench", Status: core.Running})
	if err != nil {
		t.Fatalf("ChangeWorkbenchStatus returned error: %v", err)
	}
	if out.Message != "Workbench StoppedWorkbench is running" {
		t.Errorf("expected StoppedWorkbench is running, got: %q", out.Message)
	}

	_, out, err = tools.ChangeWorkbenchStatus(context.Background(), nil, core.ChangeWorkbenchStatusInput{Namespace: "ns1", WorkbenchName: "RunningWorkbench", Status: core.Stopped})
	if err != nil {
		t.Fatalf("ChangeWorkbenchStatus returned error: %v", err)
	}
	if out.Message != "Workbench RunningWorkbench is stopped" {
		t.Errorf("expected RunningWorkbench is stopped, got: %q", out.Message)
	}
}

// TODO
func TestCreateWorkbench(t *testing.T) {
}

func newUnstructuredImage(name, displayName, repoURL string, versions []string) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(core.ImagesGVR.GroupVersion().WithKind("ImageStream"))
	u.SetName(name)
	u.SetNamespace("redhat-ods-applications")
	u.SetLabels(map[string]string{
		"opendatahub.io/notebook-image": "true",
	})
	u.SetAnnotations(map[string]string{
		"opendatahub.io/notebook-image-name": displayName,
	})

	if repoURL != "" {
		unstructured.SetNestedField(u.Object, repoURL, "status", "dockerImageRepository")
	}

	tags := make([]interface{}, len(versions))
	for i, v := range versions {
		tags[i] = map[string]interface{}{"name": v}
	}
	unstructured.SetNestedSlice(u.Object, tags, "spec", "tags")

	return u
}

func TestListImages(t *testing.T) {
	orig := resources.GetDynamicClient
	defer func() { resources.GetDynamicClient = orig }()

	scheme := runtime.NewScheme()

	image1 := newUnstructuredImage("img1", "PyTorch", "quay.io/modh/pytorch", []string{"v1", "v2"})
	image2 := newUnstructuredImage("img2", "TensorFlow", "quay.io/modh/tensorflow", []string{"latest"})

	client := dynamicfake.NewSimpleDynamicClient(scheme, image1, image2)

	resources.GetDynamicClient = func() (dynamic.Interface, error) {
		return client, nil
	}

	_, out, err := tools.ListImages(context.Background(), nil, core.ListWorkbenchesInput{})
	if err != nil {
		t.Fatalf("ListImages returned error: %v", err)
	}

	// Check PyTorch
	if !strings.Contains(out.Images, "Image: PyTorch") {
		t.Errorf("expected PyTorch in output, got: %q", out.Images)
	}
	if !strings.Contains(out.Images, "URL: quay.io/modh/pytorch") {
		t.Errorf("expected PyTorch URL in output, got: %q", out.Images)
	}
	if !strings.Contains(out.Images, "v1") || !strings.Contains(out.Images, "v2") {
		t.Errorf("expected versions v1 and v2 for PyTorch, got: %q", out.Images)
	}

	// Check TensorFlow
	if !strings.Contains(out.Images, "Image: TensorFlow") {
		t.Errorf("expected TensorFlow in output, got: %q", out.Images)
	}
}
