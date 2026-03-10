package core

type ListResourceConsumptionOutput struct {
	CPUUsage    string `json:"cpuUsage" jsonschema_description:"the CPU usage"`
	MemoryUsage string `json:"memoryUsage" jsonschema_description:"the memory usage"`
	DiskUsage   string `json:"diskUsage" jsonschema_description:"the disk usage"`
	GPUUsage    string `json:"gpuUsage" jsonschema_description:"the GPU usage"`
	UpTime      string `json:"upTime" jsonschema_description:"the up time"`
}

type ListResourceConsumptionPerWorkbenchInput struct {
	Namespace     string `json:"namespace" jsonschema_description:"the namespace of the workbench"`
	WorkbenchName string `json:"workbenchName" jsonschema_description:"the name of the workbench"`
}

type ListResourceConsumptionPerNamespaceInput struct {
	Namespace string `json:"namespace" jsonschema_description:"the namespace of the namespace"`
}

type ListResourceConsumptionPerUserInput struct {
	User string `json:"user" jsonschema_description:"the user of the user"`
}
