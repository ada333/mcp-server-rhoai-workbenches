package tools

import (
	"context"
	"strings"
	"testing"

	core "github.com/amaly/mcp-server-rhoai/core"
	"github.com/amaly/mcp-server-rhoai/resources"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

func newUnstructuredImageForToolTest(name, displayName, repoURL string, versions []string) *unstructured.Unstructured {
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

	image1 := newUnstructuredImageForToolTest("img1", "PyTorch", "quay.io/modh/pytorch", []string{"v1", "v2"})
	image2 := newUnstructuredImageForToolTest("img2", "TensorFlow", "quay.io/modh/tensorflow", []string{"latest"})

	client := dynamicfake.NewSimpleDynamicClient(scheme, image1, image2)

	resources.GetDynamicClient = func() (dynamic.Interface, error) {
		return client, nil
	}

	_, out, err := ListImages(context.Background(), nil, core.ListWorkbenchesInput{})
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
