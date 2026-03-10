package core

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

type WorkbenchInfo struct {
	Name             string `json:"name" jsonschema_description:"the name of the workbench"`
	User             string `json:"user" jsonschema_description:"the user of the workbench"`
	Status           string `json:"status" jsonschema_description:"the status of the workbench"`
	ImageDisplayName string `json:"image" jsonschema_description:"the image of the workbench"`
	ImageTag         string `json:"imageTag" jsonschema_description:"the image tag of the workbench"`
	HardwareProfile  string `json:"hardwareProfile" jsonschema_description:"the name of the hardware profile of the workbench"`
	PVCName          string `json:"pvcName" jsonschema_description:"the name of the PVC of the workbench"`
	Namespace        string `json:"namespace" jsonschema_description:"the namespace of the workbench"`
	Uptime           string `json:"uptime" jsonschema_description:"the uptime of the workbench"`
	CPUUsage         string `json:"cpuUsage" jsonschema_description:"the CPU usage of the workbench"`
	MemoryUsage      string `json:"memoryUsage" jsonschema_description:"the memory usage of the workbench"`
	DiskUsage        string `json:"diskUsage" jsonschema_description:"the disk usage of the workbench"`
	GPUUsage         string `json:"gpuUsage" jsonschema_description:"the GPU usage of the workbench"`
}

type ListWorkbenchesResult struct {
	Workbenches []WorkbenchInfo `json:"workbenches" jsonschema_description:"the list of workbenches"`
}

type ListWorkbenchesInput struct {
	Namespace string `json:"namespace" jsonschema_description:"the namespace of the workbench"`
}

type ChangeWorkbenchStatusInput struct {
	Namespace     string          `json:"namespace" jsonschema_description:"the namespace of the workbench"`
	WorkbenchName string          `json:"workbenchName" jsonschema_description:"the name of the workbench"`
	Status        WorkbenchStatus `json:"status" jsonschema_description:"the status of the workbench"`
}

type CreateWorkbenchInput struct {
	Namespace        string          `json:"namespace" jsonschema_description:"the namespace of the workbench"`
	WorkbenchName    string          `json:"workbenchName" jsonschema_description:"the name of the workbench"`
	ImageDisplayName string          `json:"imageDisplayName" jsonschema_description:"the image display name - f.e. Jupyter | Data Science | CPU | Python 3.12"`
	ImageTag         string          `json:"imageTag" jsonschema_description:"the image tag "`
	HardwareProfile  HardwareProfile `json:"hardwareProfile" jsonschema_description:"the hardware profile to use"`
	PVCName          string          `json:"pvcName" jsonschema_description:"the name of the PVC"`
}

type UpdateWorkbenchInput struct {
	Namespace        string          `json:"namespace" jsonschema_description:"the namespace of the workbench"`
	WorkbenchName    string          `json:"workbenchName" jsonschema_description:"the name of the workbench"`
	ImageDisplayName string          `json:"imageDisplayName" jsonschema_description:"the image display name - f.e. Jupyter | Data Science | CPU | Python 3.12"`
	ImageTag         string          `json:"imageTag" jsonschema_description:"the image tag "`
	HardwareProfile  HardwareProfile `json:"hardwareProfile" jsonschema_description:"the hardware profile to use"`
	PVCName          string          `json:"pvcName" jsonschema_description:"the name of the PVC"`
}

type DeleteWorkbenchInput struct {
	Namespace     string `json:"namespace" jsonschema_description:"the namespace of the workbench"`
	WorkbenchName string `json:"workbenchName" jsonschema_description:"the name of the workbench"`
}
