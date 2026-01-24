package core

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

var WorkbenchesGVR = schema.GroupVersionResource{Group: "kubeflow.org", Version: "v1", Resource: "notebooks"}

var PvcGVR = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "persistentvolumeclaims"}

var ImagesGVR = schema.GroupVersionResource{Group: "image.openshift.io", Version: "v1", Resource: "imagestreams"}

var HardwareProfilesGVR = schema.GroupVersionResource{Group: "infrastructure.opendatahub.io", Version: "v1", Resource: "hardwareprofiles"}

func GetDefaultNamespace() string {
	if ns := os.Getenv("DEFAULT_NAMESPACE"); ns != "" {
		return ns
	}
	return "redhat-ods-applications"
}

type PodsOutput struct {
	Pods string `json:"pods" jsonschema_description:"the list of pods"`
}

type ListWorkbenchesResult struct {
	Workbenches string `json:"workbenches" jsonschema_description:"the list of workbenches"`
}

type ListWorkbenchesInput struct {
	Namespace string `json:"namespace" jsonschema_description:"the namespace of the workbench"`
}

type ChangeWorkbenchStatusInput struct {
	Namespace     string          `json:"namespace" jsonschema_description:"the namespace of the workbench"`
	WorkbenchName string          `json:"workbenchName" jsonschema_description:"the name of the workbench"`
	Status        WorkbenchStatus `json:"status" jsonschema_description:"the status of the workbench"`
}

type DefaultToolOutput struct {
	Message string `json:"message" jsonschema_description:"the message with result of the tool execution"`
}

type WorkbenchStatus int

const (
	Running WorkbenchStatus = iota
	Stopped
)

// used for printing the status
func (s WorkbenchStatus) String() string {
	switch s {
	case Running:
		return "running"
	case Stopped:
		return "stopped"
	default:
		return "unknown"
	}
}

type CreateWorkbenchInput struct {
	Namespace        string          `json:"namespace" jsonschema_description:"the namespace of the workbench"`
	WorkbenchName    string          `json:"workbenchName" jsonschema_description:"the name of the workbench"`
	ImageDisplayName string          `json:"imageDisplayName" jsonschema_description:"the image display name - f.e. Jupyter | Data Science | CPU | Python 3.12"`
	ImageTag         string          `json:"imageTag" jsonschema_description:"the image tag "`
	HardwareProfile  HardwareProfile `json:"hardwareProfile" jsonschema_description:"the hardware profile to use"`
}

type ListImagesOutput struct {
	Images string `json:"images" jsonschema_description:"the list of images"`
}

type ListNamespacesOutput struct {
	Namespaces string `json:"namespaces" jsonschema_description:"the list of namespaces"`
}

type CreateCustomImageInput struct {
	ImageLocation    string `json:"imageLocation" jsonschema_description:"the location of the image"`
	ImageName        string `json:"imageName" jsonschema_description:"the name of the image"`
	ImageDescription string `json:"imageDescription" jsonschema_description:"the description of the image"`
}

type DeleteWorkbenchInput struct {
	Namespace     string `json:"namespace" jsonschema_description:"the namespace of the workbench"`
	WorkbenchName string `json:"workbenchName" jsonschema_description:"the name of the workbench"`
}

type DeleteImageInput struct {
	ImageName string `json:"imageName" jsonschema_description:"the name of the image"`
}

type HardwareProfileOutput struct {
	Message         string `json:"message" jsonschema_description:"the message with result of hardware profile creation"`
	HardwareProfile string `json:"hardwareProfile" jsonschema_description:"the hardware profile created"`
}

type HardwareProfile struct {
	HardwareProfileName string                    `json:"hardwareProfileName" jsonschema_description:"the name of the hardware profile"`
	Resources           []HardwareProfileResource `json:"resources" jsonschema_description:"the resources of the hardware profile"`
}

type HardwareProfileResource struct {
	ResourceName       string `json:"resourceName" jsonschema_description:"the name of the resource"`
	ResourceIdentifier string `json:"resourceIdentifier" jsonschema_description:"the identifier of the resource"`
	ResourceType       string `json:"resourceType" jsonschema_description:"the type of the resource"`
	DefaultCount       string `json:"defaultCount" jsonschema_description:"the default count of the resource"`
	MaxCount           string `json:"maxCount" jsonschema_description:"the max count of the resource"`
	MinCount           string `json:"minCount" jsonschema_description:"the min count of the resource"`
}

type DeleteHardwareProfileInput struct {
	HardwareProfileName string `json:"hardwareProfileName" jsonschema_description:"the name of the hardware profile"`
}
