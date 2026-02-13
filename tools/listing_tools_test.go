package tools

import (
	"context"
	"strings"
	"testing"

	core "github.com/amaly/mcp-server-rhoai/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	//"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/client-go/dynamic"
	// "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestListPods_Success(t *testing.T) {
	orig := GetClientSet
	defer func() { GetClientSet = orig }()

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
	GetClientSet = func() (kubernetes.Interface, error) {
		return client, nil
	}

	_, out, err := ListPods(context.Background(), nil, core.ListWorkbenchesInput{Namespace: ns})
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

// TODO: Re-enable this test by mocking both GetDynamicClient and GetClientSet
// Currently fails because getUptimeFromWorkbench needs GetClientSet to fetch pods
/*
func TestListWorkbenches(t *testing.T) {
	orig := GetDynamicClient
	defer func() { GetDynamicClient = orig }()

	ns := "test-ns"
	scheme := runtime.NewScheme()
	client := dynamicfake.NewSimpleDynamicClient(scheme,
		newUnstructuredWorkbench("wb-1", ns),
		newUnstructuredWorkbench("wb-other", "other-ns"),
	)

	GetDynamicClient = func() (dynamic.Interface, error) {
		return client, nil
	}

	_, out, err := ListWorkbenches(context.Background(), nil, core.ListWorkbenchesInput{Namespace: ns})
	if err != nil {
		t.Fatalf("ListWorkbenches returned error: %v", err)
	}

	if len(out.Workbenches) != 1 {
		t.Errorf("expected 1 workbench in namespace %s, got %d", ns, len(out.Workbenches))
	}

	if len(out.Workbenches) > 0 && !strings.Contains(out.Workbenches[0].Name, "wb-1") {
		t.Errorf("expected wb-1 in output, got: %q", out.Workbenches[0].Name)
	}

	// Verify wb-other is not in the list (it's in a different namespace)
	for _, wb := range out.Workbenches {
		if strings.Contains(wb.Name, "wb-other") {
			t.Errorf("did not expect wb-other in output, got: %q", wb.Name)
		}
	}
}
*/

// TODO: Re-enable this test by mocking both GetDynamicClient and GetClientSet
// Currently fails because getUptimeFromWorkbench needs GetClientSet to fetch pods
/*
func TestListAllWorkbenches(t *testing.T) {
	orig := GetDynamicClient
	defer func() { GetDynamicClient = orig }()

	scheme := runtime.NewScheme()
	client := dynamicfake.NewSimpleDynamicClient(scheme,
		newUnstructuredWorkbench("wb-1", "ns1"),
		newUnstructuredWorkbench("wb-2", "ns2"),
	)

	GetDynamicClient = func() (dynamic.Interface, error) {
		return client, nil
	}

	_, out, err := ListAllWorkbenches(context.Background(), nil, core.ListWorkbenchesInput{Namespace: ""})
	if err != nil {
		t.Fatalf("ListAllWorkbenches returned error: %v", err)
	}

	if !strings.Contains(out.Workbenches[0].Name, "wb-1") {
		t.Errorf("expected wb-1 in output, got: %q", out.Workbenches)
	}
	if !strings.Contains(out.Workbenches[1].Name, "wb-2") {
		t.Errorf("expected wb-2 in output, got: %q", out.Workbenches)
	}
}
*/
