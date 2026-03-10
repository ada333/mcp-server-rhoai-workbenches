package core

type PodsOutput struct {
	Pods string `json:"pods" jsonschema_description:"the list of pods"`
}

type DefaultToolOutput struct {
	Message string `json:"message" jsonschema_description:"the message with result of the tool execution"`
}

type ListNamespacesOutput struct {
	Namespaces string `json:"namespaces" jsonschema_description:"the list of namespaces"`
}
