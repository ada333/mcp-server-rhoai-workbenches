package core

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

var WorkbenchesGVR = schema.GroupVersionResource{Group: "kubeflow.org", Version: "v1", Resource: "notebooks"}

var PvcGVR = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "persistentvolumeclaims"}

var ImagesGVR = schema.GroupVersionResource{Group: "image.openshift.io", Version: "v1", Resource: "imagestreams"}

var HardwareProfilesGVR = schema.GroupVersionResource{Group: "infrastructure.opendatahub.io", Version: "v1", Resource: "hardwareprofiles"}

var PodsGVR = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

func GetDefaultNamespace() string {
	if ns := os.Getenv("DEFAULT_NAMESPACE"); ns != "" {
		return ns
	}
	return "redhat-ods-applications"
}
